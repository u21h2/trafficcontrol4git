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

import type { CreatedData } from "cypress/support/testing.data";

describe("Role creation/edit page", () => {
	beforeEach(() => {
		cy.login();
	});
	it("Edits an existing Role", () => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				const {role} = data;
				cy.visit(`/core/roles/${role.name}`);
				cy.get("input[name=description]").should("be.enabled").and("have.value", role.description);
				cy.get("input[name=name]").should("be.enabled").and("have.value", role.name);
				cy.get("textarea[name=permissions]")
					.should("have.value", (role.permissions ?? [""]).join("\n"))
					.and("not.have.attr", "readonly");
				cy.get("button").contains("Save").should("not.be.disabled");
			}
		);
	});

	it("Creates new Roles", () => {
		cy.visit("/core/new-role");
		cy.get("input[name=description]").should("be.enabled").and("be.empty");
		cy.get("input[name=name]").should("be.enabled").and("be.empty");
		cy.get("textarea[name=permissions]").should("be.enabled").and("not.have.attr", "readonly");
		cy.get("button").contains("Save").should("not.be.disabled");
	});
});
