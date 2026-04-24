# vhs tape for the time-tracker demo GIF.
#
# Render with:  vhs docs/demo.tape
# Output:       docs/demo.gif
#
# Requires https://github.com/charmbracelet/vhs (`brew install vhs`).

Output docs/demo.gif

Set Shell      "zsh"
Set FontSize   16
Set Width      1100
Set Height     650
Set Padding    20
Set Theme      "Catppuccin Mocha"
Set TypingSpeed 40ms

# Build once so the demo doesn't show "go build" output.
Hide
Type "go build -o ./bin/time-tracker ./cmd/time-tracker"
Enter
Sleep 4s
Show

Type "./bin/time-tracker"
Enter
Sleep 1500ms

# --- Stopwatch -----------------------------------------------------------
Type "e"
Sleep 400ms
Type "demo: writing README"
Sleep 500ms
Enter
Sleep 400ms

Space
Sleep 3s
Space
Sleep 800ms
Space
Sleep 2s
Type "s"
Sleep 1500ms

# --- Pomodoro ------------------------------------------------------------
Tab
Sleep 1s
Type "b"
Sleep 600ms
Type "+"
Sleep 400ms
Type "-"
Sleep 600ms
Space
Sleep 2500ms
Space
Sleep 800ms
Type "r"
Sleep 1s

# --- Sessions / history with filter --------------------------------------
Tab
Sleep 1500ms
Type "f"
Sleep 500ms
Type "demo"
Sleep 1500ms
Enter
Sleep 1500ms
Type "c"
Sleep 1500ms

Type "q"
Sleep 500ms
