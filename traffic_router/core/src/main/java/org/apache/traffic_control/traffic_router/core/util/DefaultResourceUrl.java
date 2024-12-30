/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.traffic_control.traffic_router.core.util;

class DefaultResourceUrl implements ResourceUrl {
	private final String[] urla;
	private int i = 0;

	public DefaultResourceUrl(final String[] urla) {
		this.urla = urla;
	}

	@Override
	public String nextUrl() {
		i++;
		i %= urla.length;
		return urla[i];
	}
}
