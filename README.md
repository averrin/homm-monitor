#HoMM Monitor

## What is it?

This is a software for gathering data from your HoMM3 game and showing it in OBS overlay. There are three parts:

- Plugin. Its a code executes inside the game to access your data. Yes, technically, its cheat, but we will make sure not to break the original game flow.
- Server. It's the central part of this program. Here we store your data, calculate analytics, and serving it for OBS.
- Widgets. When you want to show some data in OBS you need to add some source. And you may want to display different data in different places. So, there are a couple of predifined "widgets" to make it possible. Widget is a simple html page that can be displayed by "Browser" source in OBS.

## How to use it?

Go to [Releases page](https://github.com/averrin/homm-monitor/releases) and download latest.

## For commentator:

### Start server

Just double click on "HoMM-Monitor\comentator.exe"

### Create match

Type a human-readable name for your match and press the "Create Match" button. Then copy appeared Match ID and send it to your players.

### Do some tweaks

You can control some HUD settings at the "Overlay settings" tab. There are only two options for now. Also, you can see timestamps for players' updates at the "Status" tab.

### Setup OBS overlay

You can use the example "Comentator's" scene from "HoMM-Monitor\obs". Browser sources in OBS are not very handy but pretty functional.

There are three main examples:

- Commentator HUD (http://localhost:8988/widgets/commentator_hud.html) which is pretty similar to SNG Online HUD
- Map ([http://localhost:8988/widgets/map.html](http://localhost:8988/widgets/map.html))
- Heroes list ([http://localhost:8988/widgets/commentator_heroes.html](http://localhost:8988/widgets/commentator_heroes.html))

## For players:

### Install the plugin.

Unfortunately, there is no convenient way to load the plugin into HotA version. So, let's do some magic.

1. Go to folder "%your_heroes_path%\_HD3_Data\Common"
2. Backup file "cursors.dll"
3. Copy new from folder "HoMM-Monitor\plugin"

### Start server

Just double click on "HoMM-Monitor\player.exe"

### Connect to match

Got match ID from your commentator and place it to the input on "Match" tab. Also you should specify your name for this match (it will display on commentator's HUD) and press "Connect".

### Play

Preparations are done!
Now you can start the game. After map generation, even before any actions, you should press the F6 key. If everything is okay, you will see "Stats reporter started" text on the screen. Then do some actions (e.g., hero movement) and observe changes in OBS overlays. The next F6 press will stop data reporting and reset stored data. The plugin can detect restarts and will reset the game state on the server, but only if you didn't change your color. In this case, you should do it manually by pressing F6.

### Optional: Setup OBS overlay (if you play without commentator)

You can use the example "Player's" scene from "HoMM-Monitor\obs". Browser sources in OBS are not very handy but pretty functional.

There are two main examples:

- HUD (http://localhost:8989/widgets/hud.html) which is left half of the Commentator's HUD
- "Single" (http://localhost:8989/widgets/single.html?key=totalMPSpent) for displaying a single value. You can choose the desired value by changing the "key" parameter in the URL. You can take key names from the server window after data update.
- Heroes list ([http://localhost:8989/widgets/heroes.html](http://localhost:8988/widgets/commentator_heroes.html))

## Something went wrong

There are three possible ways to fix "something" without restarting everything.

- Update OBS sources. Widgets can lose connection to the server (at least it happens every time when you restart server), so you should update it. It can be done by button from the image below. If this checkbox is set, you can refresh all widgets by going to the other scene and back.

![https://s3-us-west-2.amazonaws.com/secure.notion-static.com/9cf459cf-372d-4bf6-b3c9-c6a68c6f125d/Untitled.png](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/9cf459cf-372d-4bf6-b3c9-c6a68c6f125d/Untitled.png)

- Press "Reset" button at the server GUI. It clears all accumulated analytics values (like totalMPSpent)
- Press F6 in the game to stop plugin and press again to restart it. It also resets server (see above) and reset inner clicks counter.