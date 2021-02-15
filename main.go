package main

import (
  "database/sql"
  "fmt"
  _ "github.com/mattn/go-sqlite3"
  "log"
  //"os"
  "bufio"
  "strings"
  "net"
  //"io"
)

// global listening address that people will use the telnet command to hook up to
const listenAddress = ":3410"

var allCommands = make(map[string]func(string, string, *Player))
var zones = make(map[int]*Zone)
var rooms = make(map[int]*Room)
var directions = make(map[string]int)
var player = Player{}
var db *sql.DB
var players []*Player

type Zone struct {
    ID    int
    Name  string
    Rooms []*Room
    Players []*Player
}

type Room struct {
    ID          int
    Zone        *Zone
    Name        string
    Description string
    Exits       [6]Exit
    Players []*Player
}

type Exit struct {
    To          *Room
    Description string
    Direction string
}

type Player struct {
  Room *Room
  Connection net.Conn
  Name string
}

// From notes in class 2/8

type InputEvent struct {
  Player *Player
  //Connection net.Conn
  Command string
  Close bool
  Login bool
}

type OutputEvent struct {
  Text string
}

/*
// player interface
func (p *Player) Printf(format string, a ...interface{}) {
  msg := fmt.Sprintf(format, a...)
  // being sent to the Output job, that's the difference between = and <-
  p.Output <- OutputEvent(Text: msg)
}

func handleConnection(db *sql.DB, conn net.Conn, inputs chan<- InputEvent) {
  defer func() {

  }
}
type Renderer interface{}
*/
/*
func commandLoop() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		doCommand(line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("in main command loop: %v", err)
	}

	return nil
}
*/

func addCommand(cmd string, f func(string, string, *Player)) {
	for i := range cmd {
		if i == 0 {
			continue
		}
		prefix := cmd[:i]
		allCommands[prefix] = f
	}
	allCommands[cmd] = f
}

func initCommands() {
	addCommand("smile", cmdSmile)
	addCommand("south", cmdSouth)
  addCommand("north", cmdNorth)
  addCommand("east", cmdEast)
  addCommand("west", cmdWest)
  addCommand("tiphat", cmdTipHat)
  addCommand("look", cmdLook)
  addCommand("up", cmdUp)
  addCommand("down", cmdDown)
  addCommand("recall", cmdRecall)
  addCommand("whisper", cmdWhisper)
  addCommand("say", cmdSay)
  addCommand("shout", cmdShout)
  addCommand("gossip", cmdGossip)
  addCommand("quit", cmdQuit)
}

func is_prefix(s, cmd string) bool {
  for i := range cmd {
    if s == cmd[i:] {
      return true
    }
  }
  return false
}

func doCommand(cmd string, player *Player) error {
  words := strings.Fields(cmd)
	if len(words) == 0 {
		return nil
	} else if len(words) == 1 {
    if f, exists := allCommands[strings.ToLower(words[0])]; exists {
      f("", "", player)
        return nil
    }
  } else if len(words) >= 2 {
		if is_prefix(words[0], "look"){
      // ToLower for userproof because users are ducking idiots
			if f, exists := allCommands[strings.ToLower(words[0])]; exists {
				f(words[1], "", player)
				return nil
			}
		} else if is_prefix(words[0], "say") || is_prefix(words[0], "shout") || is_prefix(words[0], "gossip") {
      if f, exists := allCommands[strings.ToLower(words[0])]; exists {
        f(strings.Join(words[1:], ""), "", player)
          return nil
      }
    } else if is_prefix(words[0], "whisper") {
      if _, exists := allCommands[strings.ToLower(words[0])]; exists {
        cmdWhisper(words[1], strings.Join(words[2:], " "), player)
        return nil
      }
    }
	}
  fmt.Fprintf(player.Connection, words[0] + " is not a command\n")
	return nil
}

// directional commands

func cmdNorth(s string, s0 string, player *Player) {
  if len(player.Room.Exits[0].Description) == 0 {
    fmt.Fprintf(player.Connection, "You can't go north from here, you're like an unatractive security guard giving me a patdown at an airport. \n")
  } else {
    // .To is a pointer to another room, man I hate go
    player.Room = player.Room.Exits[0].To
    fmt.Fprintf(player.Connection, "You proceed north.")
    fmt.Fprintf(player.Connection, player.Room.Description)
    fmt.Fprintf(player.Connection, "\n")

  }
}

func cmdSouth(s string, s0 string, player *Player) {
  if len(player.Room.Exits[2].Description) == 0 {
    fmt.Fprintf(player.Connection, "You can't go south from here, I bet I could write a book about how much you don't know... \n")
  } else {
    // .To is a pointer to another room, man I hate go
    player.Room = player.Room.Exits[2].To
    fmt.Fprintf(player.Connection, "You proceed south.")
    fmt.Fprintf(player.Connection, player.Room.Description)
    fmt.Fprintf(player.Connection, "\n")

  }
}

func cmdEast(s string, s0 string, player *Player) {
  if len(player.Room.Exits[1].Description) == 0 {
    fmt.Fprintf(player.Connection, "You can't go east from here, your wife can't cheat on you because she set her standards so low marrying you... \n")
  } else {
    // .To is a pointer to another room, man I hate go
    player.Room = player.Room.Exits[1].To
    fmt.Fprintf(player.Connection, "You proceed east.")
    fmt.Fprintf(player.Connection, player.Room.Description)
    fmt.Fprintf(player.Connection, "\n")

  }
}

func cmdWest(s string, s0 string, player *Player) {
  if len(player.Room.Exits[3].Description) == 0 {
    fmt.Fprintf(player.Connection,"You can't go west from here, the middle finger was invented because of you... \n")
  } else {
    // .To is a pointer to another room, man I hate go
    player.Room = player.Room.Exits[3].To
    fmt.Fprintf(player.Connection, "You proceed west.")
    fmt.Fprintf(player.Connection,player.Room.Description)
    fmt.Fprintf(player.Connection,"\n")

  }
}

func cmdUp(s string, s0 string, player *Player) {
  if len(player.Room.Exits[4].Description) == 0 {
    fmt.Fprintf(player.Connection,"You can't go up from here, I'll never forget the first time we met... But I'll keep trying, \n")
  } else {
    // .To is a pointer to another room, man I hate go
    player.Room = player.Room.Exits[4].To
    fmt.Fprintf(player.Connection,"You proceed up.")
    fmt.Fprintf(player.Connection,player.Room.Description)
    fmt.Fprintf(player.Connection,"\n")

  }
}

func cmdDown(s string, s0 string, player *Player) {
  if len(player.Room.Exits[5].Description) == 0 {
    fmt.Fprintf(player.Connection, "You can't go up from here, idiot... \n")
  } else {
    // .To is a pointer to another room, man I hate go
    player.Room = player.Room.Exits[5].To
    fmt.Fprintf(player.Connection, "You proceed down.")
    fmt.Fprintf(player.Connection, player.Room.Description)
    fmt.Fprintf(player.Connection, "\n")

  }
}

// emote commands

func cmdSmile(s string, s0 string, player *Player) {
	fmt.Fprintf(player.Connection, "You smile happily.\n")
}

func cmdTipHat(s string, s0 string, player *Player) {
  fmt.Fprintf(player.Connection, "You tip your hat politely.\n")
}

// action commands

func cmdLook(s string, s0 string, player *Player) {
  // finds player room name and description, uses that to find the exits corresponding to each room
  if len(s) == 0 {
		fmt.Fprintln(player.Connection, player.Room.Name, "\n", player.Room.Description, "\n")
		for i, exit := range player.Room.Exits  {
			if len(exit.Description) > i-i {
				fmt.Fprintln(player.Connection, exit.Direction, exit.Description, "\n")
			}
		}
  } else if s == "north" || s == "nort" || s == "nor" || s == "no" || s == "n" {
		var desc = player.Room.Exits[0].Description
		if len(desc) > 0 {
			fmt.Fprintln(player.Connection, player.Room.Exits[0].Description)
		} else {
			fmt.Fprintln(player.Connection, "There is no where to go in this direction...\n")
		}

	} else if s == "east" || s == "eas" || s == "ea" || s == "e" {
		var desc = player.Room.Exits[1].Description
		if len(desc) > 0 {
			fmt.Fprintln(player.Connection, player.Room.Exits[1].Description)
		} else {
			fmt.Fprintln(player.Connection, "There is no where to go in this direction...\n")
		}
	} else if s == "south" || s == "sout" || s == "sou" || s == "so" || s == "s" {
		var desc = player.Room.Exits[2].Description
		if len(desc) > 0 {
			fmt.Fprintln(player.Connection, player.Room.Exits[2].Description)
		} else {
			fmt.Fprintln(player.Connection, "There is no where to go in this direction...\n")
		}
	} else if s == "west" || s == "wes" || s == "we" || s == "w" {
		var desc = player.Room.Exits[3].Description
		if len(desc) > 0 {
			fmt.Fprintln(player.Connection, player.Room.Exits[3].Description)
		} else {
			fmt.Fprintln(player.Connection, "There is no where to go in this direction...\n")
		}
	} else if s == "up" || s == "u" {
		var desc = player.Room.Exits[4].Description
		if len(desc) > 0 {
			fmt.Fprintln(player.Connection, player.Room.Exits[4].Description)
		} else {
			fmt.Fprintln(player.Connection,"There is no where to go in this direction...\n")
		}
	} else if s == "down" || s == "dow" || s == "do" || s == "d" {
		var desc = player.Room.Exits[5].Description
		if len(desc) > 0 {
			fmt.Fprintln(player.Connection,player.Room.Exits[5].Description)
		} else {
			fmt.Fprintln(player.Connection,"There is no where to go in this direction...\n")
		}
	}
}

// Whisper, say, shout, take a shit, and gossip

func cmdWhisper(s string, m string, player *Player) {
  var player_exists = false
  for _, p := range players {
    if m == "" {
      player_exists = true
      fmt.Fprintln(player.Connection, "Alright buhddy, you are trying to whisper to " + p.Name + "but you can't whisper to them.\n")
      break
    }
    if p.Name == s {
      player_exists = true
      fmt.Fprintln(player.Connection, "To (" + p.Name + ") " + m, "\n")
      fmt.Fprintln(p.Connection, "From (" + player.Name + ") " + m, "\n")
    }
  }
  if player_exists == false {
    fmt.Fprintln(player.Connection, s + "does not exist or hasn't logged in.\n")
  }
}

func cmdSay(s string, s2 string, player *Player) {
  for _, p := range player.Room.Players {
    if s == "" {
      fmt.Fprintln(player.Connection, "You unsuccessfully said anything you donkey\n")
      break
    }
    fmt.Fprintln(p.Connection, "(" + player.Name + ") ", s, "\n")
  }
}

func cmdShout(s string, s2 string, player *Player) {
  for _, p := range player.Room.Zone.Players {
    if s == "" {
      fmt.Fprintln(player.Connection, "You shout, but look like if Gordan Ramsey did meth instead\n")
      break
    }
    fmt.Fprintln(p.Connection, "(" + player.Name + ") ", s, "\n")
  }
}

func cmdGossip(s string, s2 string, player *Player) {
  for _, p := range players {
    if s == "" {
      fmt.Fprintln(player.Connection, "You try to gossip, but end up being cringe\n")
      break
    }
    fmt.Fprintln(p.Connection, "(" + player.Name + ") ", s, "\n")
  }
}

// Mud Part 2

// Objectives 1 and 2

func readZones(stmt *sql.Stmt) error {
  // new iteration of readZones

  // allocate the Query
  rows, err := stmt.Query()
  // error checker
  if err != nil {
    log.Fatalf("zone query: %v", err)
  }
  defer rows.Close()
  // reads each of the zones
  for rows.Next() {
    var id int
    var name string
    var rooms []*Room
    // the actual scan itself + an error checker
    if err := rows.Scan(&id, &name); err != nil {
      log.Fatalf("reading zones: %v", err)
    }
    // creates the zone
    var zone = Zone{id, name, rooms, players}
    zones[id] = &zone
  }
  // final error checker
  if err := rows.Err(); err != nil {
    log.Fatal(err)
  }
  return nil

  /*
  // creating the path to the database
  path := "world.db"

  options :=
    "?" + "_busy_timeout=10000" +
      "&" + "_foreign_keys=ON" +
      "&" + "_journey_mode=WAL" +
      "&" + "_synchronous=NORMAL"

  // launching the command to open the file
  db, err := sql.Open("sqlite3", path + options)
  // if we have trouble opening the database we launch an error
  if err != nil {
    log.Fatalf("opening database: %v", err)
  }
  defer db.Close()

  // Launches a query
  rows, err := db.Query(`SELECT * FROM zones`)
  // if we have trouble with the query we launch an error
  if err != nil {
    log.Fatalf("zone query: %v", err)
  }
  defer rows.Close()

  // creates a map (kinda like a dictionary) of zones
  var zones = make(map[int]*Zone)
  // loops and appends to zone
  for rows.Next() {
    var id int
    var name string
    var rooms []*Room
    if err := rows.Scan(&id, &name); err != nil {
      log.Fatalf("reading zones: %v", err)
    }
    // creates the Zone object, pushing the ID, name, and room into a zone object
    var zone = Zone{id, name, rooms}
    zones[id] = &zone
    fmt.Printf("id::%d name: %s\n", id, name)
  }
  */
}

// objective 4

func readRooms(stmt *sql.Stmt) (map[int]*Room, error) {
  rows, err := stmt.Query()
  if(err != nil) {
    log.Fatalf("room query: %v", err)
  }
  defer rows.Close()

  var rooms = make(map[int]*Room)

  for rows.Next() {
    var id int
    var zone_id int
    var name string
    var description string
    var exits [6]Exit
    if err := rows.Scan(&id, &zone_id, &name, &description); err != nil {
      log.Fatalf("reading rooms: %v", err)
    }
    var room = Room{id, zones[zone_id], name, description, exits, players}
    rooms[id] = &room
    zones[zone_id].Rooms = append(zones[zone_id].Rooms, &room)
  }
  if err := rows.Err(); err != nil {
    log.Fatal(err)
  }
  return rooms, nil
}

func readExits(stmt *sql.Stmt) error {
  // open the query
  rows, err := stmt.Query()
  // error checker
  if(err != nil) {
    log.Fatalf("exit query: %v", err)
  }
  defer rows.Close()

  // loops through, runs through the scan
  for rows.Next() {
    var from_room_id int
    var to_room_id int
    var direction string
    var description string
    // scans the rows using pointers
    if err := rows.Scan(&from_room_id, &to_room_id, &direction, &description); err != nil {
      log.Fatalf("reading exits: %v", err)
    }
    // grabs the exit
    var exit = Exit{rooms[to_room_id], description, direction}
    // goes to the from room id, grabs the direction from the array of directions, then sets all that to exit
    // tldr seting the corresponding room to the corresponding exit
    rooms[from_room_id].Exits[directions[direction]] = exit
  }
  if err := rows.Err(); err != nil {
    log.Fatal(err)
  }
  player.Room = rooms[3001]
  fmt.Println(player.Room.Name, "\n", player.Room.Description)
  return nil
}

// part 3

/* old listener
func Listener() {
  // Listen on TCP port 8080 on all available unicast and anycast IP addresses of the local system
  ln, err := net.Listen("tcp", ":8080")
  if err != nil {
    log.Fatalf("listening for tcp stuffs: %v", err)
  }
  defer ln.Close()
  for {
    conn, err := ln.Accept()
    if err != nil {
      log.Fatalf("nil was passed through the tcp connection: %v", err)
    }
    // handle connection with a new goroutine.
    // the loop then return to accepting, so that
    // multiple connection may be served concurrently
    go func(c net.Conn) {
      // Echo all incoming data
      io.Copy(c, c)
      // shut down the connection
    }(conn)
  }

}
*/

// listener
func listenForConnections(inputChannel chan InputEvent) {
  // Listens on TCP port 3410 on all avaliable unicasts and anycast IP addresses of the local system
  // keeping in mind listenAddress is a global variable delcared at the top of the program
  listen, err := net.Listen("tcp", listenAddress)
  if err != nil {
    log.Fatalf("listening on %s: %v", listenAddress, err)
  }
  defer listen.Close()

  // infinate loop that looks for and accepts single connections
  for {
    // err variable is a new one because it is inside the loop
    conn, err := listen.Accept()
    if err != nil {
      log.Fatalf("accept: %v", err)
    }
    // calls goroutine
    go handleConnection(conn, inputChannel)
  }

}

func handleConnection(conn net.Conn, inputChannel chan InputEvent) {
  // first param of Fprintf is an io writer
  fmt.Fprintf(conn, "Name: ")
  var name string
  // scan for name
  fmt.Fscanln(conn, &name)
  fmt.Fprintln(conn)
  // sends player to default room
  var player = Player{Room: rooms[3001], Connection: conn, Name: name}
  // adds the player
  players = append(players, &player)
  player.Room.Zone.Players = append(player.Room.Zone.Players, &player)
  player.Room.Players = append(player.Room.Players, &player)

  for _, p := range players {
    // reads off online players
    fmt.Fprintln(p.Connection, name, " is online\n")
  }

  fmt.Fprintln(conn, player.Room.Name, "\n", player.Room.Description)
  scanner := bufio.NewScanner(conn)
  for scanner.Scan() {
    line := scanner.Text()
    event := InputEvent{ Player: &player, Command: line }
    inputChannel <- event
  }
  if err := scanner.Err(); err != nil {
    fmt.Fprintf(player.Connection, "connection error: %v", err)
  } else {
    fmt.Fprintf(player.Connection, "connection closed normally")
  }


/*
  // scanner is an interface for reading data of a file, imported class from bufio
  scanner := bufio.NewScanner(conn)
  // looping through the scan token that .Scan returns, I believe will be in this canse a keyword
  for scanner.Scan() {
    line := scanner.Text()
    event := InputEvent{ Connection: conn, Command: line }
    inputChannel <- event
    fmt.Fprintf(conn, "you typed: %s\n", line)
  }
  if err := scanner.Err(); err != nil {
    log.Printf("connection error: %v", err)
  } else {
    log.Printf("connection closed normally\n")
  }
*/
}

func main() {

  // allocating directions
  directions["n"] = 0
  directions["e"] = 1
  directions["s"] = 2
  directions["w"] = 3
  directions["u"] = 4
  directions["d"] = 5

  log.SetFlags(log.Ltime | log.Lshortfile)

  // put the world.db path in
  path := "world.db"

  // allocates the options
  options :=
    "?" + "_busy_timeout=10000" +
      "&" + "_foreign_keys=ON" +
      "&" + "_journey_mode=WAL" +
      "&" + "_synchronous=NORMAL"

  // launching the command to open the file
  db, err := sql.Open("sqlite3", path + options)
  // if we have trouble opening the database we launch an error
  if err != nil {
    log.Fatalf("opening database: %v", err)
  }
  defer db.Close()

  // Transaction Section
  tx, err := db.Begin()
  if(err != nil) {
    log.Fatalf("begin zone read transaction: %v", err)
  }
  stmt, err := tx.Prepare(`SELECT * FROM zones`)
  if(err != nil) {
    log.Fatalf("prepare zone read transaction: %v", err)
  }
  defer stmt.Close()

  err = readZones(stmt)
  // commits and rollbacks
  if(err != nil) {
    // if it fails (doesn't return nil) we rollback the transaction
    tx.Rollback()
  } else {
    // nil has been returned, we shall now commit
    tx.Commit()
  }

  tx, err = db.Begin()
  if(err != nil) {
    log.Fatalf("begin room read transaction: %v", err)
  }
  stmt, err = tx.Prepare(`SELECT * FROM rooms`)
  if err != nil {
    log.Fatalf("prepare room read transaction: %v", err)
  }

  defer stmt.Close()
  rooms, err = readRooms(stmt)
  if(err != nil) {
    tx.Rollback()
  } else {
    tx.Commit()
  }

  // Exit Read transaction
  tx, err = db.Begin()
  if(err != nil) {
    log.Fatalf("begin exit read transaction: %v", err)
  }
  stmt, err = tx.Prepare(`SELECT * FROM exits`)
  if err != nil {
    log.Fatalf("prepare exit read transaction: %v", err)
  }
  defer stmt.Close()
  err = readExits(stmt)
  if(err != nil) {
    tx.Rollback()
  } else {
    tx.Commit()
  }

  initCommands()

  // part 3 (I think it goes here)

  // create the channel with all the workers
  inputChannel := make(chan InputEvent)
  go listenForConnections(inputChannel)
  mainLoop(inputChannel)
  // end of server connection thing

/*
  if err := commandLoop(); err != nil {
    log.Fatalf("%v", err)

  }
*/


}

// listening for server events
func mainLoop(inputChannel chan InputEvent) {
  for event := range inputChannel {
    doCommand(event.Command, event.Player)
  }
}

func cmdRecall(s string, s0 string, player *Player) {
  player.Room = rooms[3001]
  fmt.Fprintf(player.Connection, "You have returned to the Temple of Midgaaurd. \n")
}

func cmdQuit(s string, s0 string, player *Player) {
  fmt.Fprintf(player.Connection, "Catch you later now! \n")
  player.Connection.Close()
}

/*
From all my viewers who were kind enough to understand I couldn't stream tonight because of this assignment, I thank you.
You all are a bunch of chads.
Oh and Curtis, if you're watching, you're a chad too.

if (assignmentTurnedIn = true)
{
      jexGrade == 100;
}
*/
