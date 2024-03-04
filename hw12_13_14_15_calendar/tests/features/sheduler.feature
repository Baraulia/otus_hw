# file: scheduler.feature

Feature: sending notification about events
    As client of scheduler service
	In order to understand that the user was informed about event
	I want to receive delivery confirmation from notifications queue

	Scenario: Notification about event is received
		Given Setup test for scheduler
	    When Sender processed a new notification I want to receive notification about it
		Then Teardown test for scheduler
