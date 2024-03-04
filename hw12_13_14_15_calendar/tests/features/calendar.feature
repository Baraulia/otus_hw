# file: calendar.feature

# http://calendar:8085/

Feature: HTTP API for CRUD Operations on Database

    Scenario: Creating an event in the database
        Given Setup test for calendar
        When I send "POST" request to "http://calendar:8085/event"
        Then The response code should be 201
        And Teardown test for calendar

    Scenario: updating an event in the database
        Given Setup test for calendar
        And Create new event
        When I send "PUT" request to "http://calendar:8085/event/"
        Then The response code should be 204
        And Teardown test for calendar

    Scenario: receiving list events from database
        Given Setup test for calendar
        And Create new event
        When I send "GET" request to "http://calendar:8085/event/list"
        Then The response code should be 200
        And The response should match text in variable listEvents
        And Teardown test for calendar

    Scenario: deleting an event in the database by id
        Given Setup test for calendar
        And Create new event
        When I send "DELETE" request to "http://calendar:8085/event/"
        Then The response code should be 204
        And Teardown test for calendar