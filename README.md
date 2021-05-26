# go-mon
A file watcher system written using golang. It helps in restarting the runtime if any change is detected.
At any point you can press ```rs``` to restart the command.

# Command Flags
```
  -c = Command to re-run the system (default : make run)
  -f = Folder to watch out for (default : .)
  -i = Ignore files starting with . (default : true)
  -e = File Extensions to watch out for (default : .go)
```

# Examples

For Watching the templates folder
```
go-mon -f ./templates
```

For Watching the folder with javascript files
```
go-mon -e .js,.html
```
