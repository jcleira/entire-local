# entire-local Actions Demo
# Run: vhs demo.tape
# Requires: entire-local binary built, repo with entire/checkpoints/v1 branch

Output demo.gif

Set Shell "bash"
Set FontSize 15
Set FontFamily "Noto Sans Mono"
Set Width 1400
Set Height 800
Set Theme "Dracula"
Set WindowBar Colorful
Set BorderRadius 8
Set Padding 10
Set TypingSpeed 0
Set CursorBlink false

# Launch TUI
Type "./entire-local"
Enter
Sleep 2s

# Show help screen with new 'a' binding
Type "?"
Sleep 2s
Type "?"
Sleep 500ms

# Open action menu on overview (checkpoint-specific items dimmed)
Type "a"
Sleep 2s

# Navigate through action items
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms

# Navigate back up
Type "k"
Sleep 400ms
Type "k"
Sleep 400ms
Type "k"
Sleep 1s

# Close action menu
Escape
Sleep 1s

# Navigate to a checkpoint and enter detail view
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Enter
Sleep 2s

# Open action menu on detail (all items enabled)
Type "a"
Sleep 2s

# Select rewind (destructive) to show confirmation dialog
Type "w"
Sleep 2s

# Cancel the confirmation
Type "n"
Sleep 1s

# Open action menu again and try status
Type "a"
Sleep 1s
Type "s"
Sleep 3s

# Close status output
Escape
Sleep 1s

# Quit
Type "q"
Sleep 500ms
