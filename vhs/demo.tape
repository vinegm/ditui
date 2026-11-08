# A seeded drill for the README: decoding with the cheatsheet, then the breakdowns

Output assets/demo.gif

Set Shell bash
Set FontFamily "JetBrainsMono Nerd Font, JetBrains Mono, monospace"
Set FontSize 18
Set Width 900
Set Height 700
Set Padding 20
Set TypingSpeed 100ms
Set Theme "TokyoNight"

# Do setup invisibly
Hide
Type "./ditui --config Themes/tokyo-night.yaml --seed demo"
Enter
Sleep 500ms
Show

# Start capturing
Sleep 1s

# Turn on the cheatsheet
Type@300ms "jjj"
Sleep 300ms
Enter
Sleep 800ms

# Start the drill
Type@200ms "kkk"
Sleep 300ms
Enter
Sleep 1.5s

# Answer a few, mixing up Y and K on the third
Type "x"
Sleep 1.2s
Type "w"
Sleep 1.5s
Type "k"
Sleep 2.5s
Type "s"
Sleep 1s
Type "b"
Sleep 1.2s
Type "v"
Sleep 1.2s
Type "q"
Sleep 1.5s

# Finish, and go through the breakdowns
Escape
Sleep 2.5s
Type "l"
Sleep 2s
Type "l"
Sleep 3s
