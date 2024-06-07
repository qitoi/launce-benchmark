from locust import User, task, constant

import random

exc = Exception("unexpected response status code")

class MyUser(User):
    @task
    def foo(self):
        response_time = random.random() * 1000
        content_length = random.randrange(1024 * 1024)
        self.environment.events.request.fire(
            request_type="GET",
            name="/foo",
            response_time=response_time,
            response_length=content_length,
        )

    @task
    def bar(self):
        response_time = random.random() * 1000
        content_length = random.randrange(1024 * 1024)
        self.environment.events.request.fire(
            request_type="GET",
            name="/bar",
            response_time=response_time,
            response_length=content_length,
            exception=exc,
        )
