# ChessEngine

The goal of my project is to create a chess engine that is built generically in order to work for variations of the chess ruleset and eventually other turn-based board games.  So far, I have built a basic engine that can play normal chess slowly and a UI to test it.


Update:
Changed the data structure of a Board from an array to map
Implemented some mutations but sometimes it is necessary to create new boards
This is because some new boards are used to calculate hypothetical moves
And should not affect the original board
