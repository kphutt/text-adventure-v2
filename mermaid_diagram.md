```mermaid
graph TD
    %% Define main actors and systems
    subgraph User
        Player[Player]
    end

    subgraph "Application: text-adventure-v2"
        direction TB
        MainApp["main.go<br><b>(View/Controller)</b>"]
        GameCore["Game Package<br><b>(Model)</b>"]
    end

    subgraph "Game Package Details"
        direction LR
        GameStruct[<b>game.go</b><br>Game State & Logic]
        Structs[<b>structs.go</b><br>Data Structures]
        World[<b>world.go</b><br>World Definition]
        Parser[<b>parser.go</b><br>Input Parser]
        MapGen[<b>map.go</b><br>Map Generator]
    end

    %% --- Connections ---
    Player -- "Types commands" --> MainApp
    MainApp -- "Initializes & Interacts" --> GameCore
    GameCore -- "Comprises" --> GameStruct
    GameCore -- "Comprises" --> Structs
    GameCore -- "Comprises" --> World
    GameCore -- "Comprises" --> Parser
    GameCore -- "Comprises" --> MapGen

    MainApp -- "Calls methods on" --> GameStruct
    GameStruct -- "Returns output" --> MainApp

    GameStruct --> Structs
    GameStruct --> World
    GameStruct --> Parser
    GameStruct --> MapGen

    %% --- Styling for clarity ---
    classDef actor fill:#ADD8E6,stroke:#000,stroke-width:2px,font-weight:bold;
    classDef container fill:#FFFF99,stroke:#333,stroke-width:2px,font-weight:bold;
    classDef package fill:#B0E0E6,stroke:#333,stroke-width:2px,font-weight:bold;
    classDef component fill:#E6E6FA,stroke:#333,stroke-width:1px,font-weight:bold;

    class Player actor;
    class MainApp container;
    class GameCore package;
    class GameStruct,Structs,World,Parser,MapGen component;

    linkStyle 0 stroke-width:2px,fill:none,stroke:#000;
    linkStyle 1 stroke-width:2px,fill:none,stroke:#000;
    linkStyle 2 stroke-width:2px,fill:none,stroke:#000;
    linkStyle 3 stroke-width:2px,fill:none,stroke:#000;
    linkStyle 4 stroke-width:2px,fill:none,stroke:#000;
    linkStyle 5 stroke-width:2px,fill:none,stroke:#000;
    linkStyle 6 stroke-width:2px,fill:none,stroke:#000;
    linkStyle 7 stroke-width:2px,fill:none,stroke:#000;
    linkStyle 8 stroke-width:2px,fill:none,stroke:#000;
    linkStyle 9 stroke-width:2px,fill:none,stroke:#000;
    linkStyle 10 stroke-width:2px,fill:none,stroke:#000;
```