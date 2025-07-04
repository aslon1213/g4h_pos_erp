from locust import HttpUser, task, between
import os


class AuthUser(HttpUser):
    wait_time = between(1, 3)
    branch_id = "2dbbbb13-60be-425e-b07c-fa4889f0400d"
    journal_id = "649e78a656b78aefd50372e4"
    supplier_id = "f5602699-8487-4a44-8ef2-eec9815dfaf8"
    transaction_id = "68656b715dfc7f1ed81b8bad"
    host = os.getenv("HOST")
    if host == "":
        print("HOST is not set")
        exit(1)

    def on_start(self):
        """This runs when a virtual user starts"""
        # Example login payload
        payload = {"username": "aslon", "password": "aslon"}

        # Login request
        with self.client.post(
            "/auth/login", json=payload, catch_response=True
        ) as response:
            if response.status_code == 200:
                # Extract token from response (e.g. JWT)
                token = response.json().get("data")
                if token:
                    # Save headers with token for future use
                    self.client.headers.update(
                        {
                            "Authorization": f"{token}",
                            "Content-Type": "application/json",
                        }
                    )
                else:
                    response.failure("No token in response")
            else:
                response.failure(f"Login failed: {response.status_code}")

    @task
    def get_journals(self):
        """Authenticated GET request"""
        resp = self.client.get(f"/api/journals/branch/{self.branch_id}")
        if resp.status_code == 400:
            print(resp.text)

    @task
    def get_journal(self):
        resp = self.client.get(f"/api/journals/{self.journal_id}")
        if resp.status_code == 400:
            print(resp.text)

    @task
    def get_all_suppliers(self):
        resp = self.client.get("/api/suppliers")
        if resp.status_code == 400:
            print(resp.text)

    @task
    def get_supplier(self):
        resp = self.client.get(f"/api/suppliers/{self.supplier_id}")
        if resp.status_code == 400:
            print(resp.text)

    @task
    def get_transactions(self):
        resp = self.client.get(f"/api/transactions/branch/{self.branch_id}")
        if resp.status_code == 400:
            print(resp.text)

    @task
    def get_transaction(self):
        resp = self.client.get(f"/api/transactions/{self.transaction_id}")
        if resp.status_code == 400:
            print(resp.text)

    # query products later
