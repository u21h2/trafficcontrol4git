#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

"""API Contract Test Case for system_info endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_system_info_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]]) -> None:
	"""
	Test step to validate keys, values and data types from system_info endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	"""
	# validate system_info keys from system_infos get response
	logger.info("Accessing /system_infos endpoint through Traffic ops session.")

	system_info_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_system_info()
	try:
		first_system_info = system_info_get_response[0]
		if not isinstance(first_system_info, dict):
			raise TypeError("malformed API response; first system_info in response is not an dict")
		logger.info("system_info Api get response %s", first_system_info)

		system_info_response_template = response_template_data.get("system_info")
		if not isinstance(system_info_response_template, dict):
			raise TypeError(
				f"system_info response template data must be a dict, not '{type(system_info_response_template)}'")

		# validate keys, data types and values from system_infos get json response.
		assert validate(instance=first_system_info, schema=system_info_response_template) is None
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for system_info endpoint: API response was malformed")
