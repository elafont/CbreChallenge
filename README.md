# Cbre Challenge

Test code for the CBRE challenge

## Go’s Hangman Game

Build a simple version of a ‘hangman’ game as a server-side application.
Requirements

1. As a user I want to be able to start a new game;
2. As a user I want to be able to guess a character for an ongoing game;
3. As a user I want to be notified when the game I am playing ends (win/game over);
4. As a user I want to resume an incomplete game;
5. As a user I want to be able to list all the games that have been played so far and the ones that are currently ongoing;
6. As a user I want to be notified when I trigger an action which results in a server failure.

## Technology

● Use Golang both for the server and client side;
● You may choose any additional technology if desired or if applicable to meet the
requirements (considering that we are a company operating on Peer to Peer networks).

## Interaction

● For the interaction with the server we expect you to provide a simple CLI which issues
commands to the server;
● Also, assume that:
  ○ No user authentication is required;
  ○ No user registration is required;
  ○ No persistence is required across restarts;

## Assessment

● Primarily, we’re looking for clean and readable code and a sound design approach.
● We will be evaluating potential concurrency issues.
● Unit Tests: To give an indication of your unit testing ability/understanding, we’d like to see example unit tests on the business logic that involve the gameplay of guessing characters. You’re welcome to add a short Go comment explaining what type of tests, coverage, etc... such-and-such could be added given more time.

## Notes

Words will be chosen from the "/usr/share/dict/british-english" dictionary file, it has been preprocessed to eliminate words with accents, or with apostrophes or upper letters words, also the preparation was done just to have a copy local to this challenge, as this dictionary may not be available on all linux systems.
