/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

describe("ASN table page", () => {
	beforeEach(() => {
		cy.login();
	});
	it("Loads elements", async () => {
		cy.visit("/core/asns");
		cy.find("input[name=fuzzControl]");
		cy.find("div.ag-row").should("have.length.at.least", 1);
	});
});
