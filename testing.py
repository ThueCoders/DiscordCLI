import sys,os
import threading
import curses
from curses import wrapper
import json
import aiofiles
import discord
import asyncio

from curses.textpad import Textbox, rectangle

class Bot(discord.Client):
    def __init__(self):
        super(Bot, self).__init__()

    async def on_ready(self):
        x = 1-1
        #print('Logged in as')
        #print(self.user.name)
        #print(self.user.id)
        #print('------')

    async def on_message(self, message):
        if message.author == self.user:
            content = message.content  # type: str

            if content.startswith('$test'):
                await self.send_message(message.channel, 'hmmmm')

def draw_menu(stdscr):
    k = 0
    cursor_x = 0
    cursor_y = 0

    # Clear and refresh the screen for a blank canvas
    stdscr.clear()
    stdscr.refresh()

    # Start colors in curses
    curses.start_color()
    curses.init_pair(1, curses.COLOR_CYAN, curses.COLOR_BLACK)
    curses.init_pair(2, curses.COLOR_RED, curses.COLOR_BLACK)
    curses.init_pair(3, curses.COLOR_BLACK, curses.COLOR_WHITE)

    # Loop where k is the last character pressed
    while (True):

        if(k == ord('q')):
            client.logout()
            break
        # Initialization
        stdscr.clear()
        height, width = stdscr.getmaxyx()

        if k == curses.KEY_DOWN:
            cursor_y = cursor_y + 1
        elif k == curses.KEY_UP:
            cursor_y = cursor_y - 1
        elif k == curses.KEY_RIGHT:
            cursor_x = cursor_x + 1
        elif k == curses.KEY_LEFT:
            cursor_x = cursor_x - 1

        cursor_x = max(0, cursor_x)
        cursor_x = min(width-1, cursor_x)

        cursor_y = max(0, cursor_y)
        cursor_y = min(height-1, cursor_y)

        # Declaration of strings
        title = "Discord CLI"[:width-1]
        subtitle = "Written by Nathan Withers"[:width-1]
        with open("ascii60Wide.txt") as f:
                content = f.readlines()
                asciiArt = [x.strip() for x in content] 
        statusbarstr = "Press 'q' to exit | STATUS BAR | Pos: {}, {}".format(cursor_x, cursor_y)

        # Centering calculations
        start_x_title = int((width // 2) - (len(title) // 2) - len(title) % 2)
        start_x_subtitle = int((width // 2) - (len(subtitle) // 2) - len(subtitle) % 2)
        start_x_art = int((width // 2) - (len(asciiArt[1]) // 2) - len(asciiArt[1]) %2)
        start_y = int((height // 2) - 2)
        start_y_art = 1

        # Rendering some text
        whstr = "Width: {}, Height: {}".format(width, height)
        stdscr.addstr(0, 0, whstr, curses.color_pair(1))

        # Render status bar
        stdscr.attron(curses.color_pair(3))
        stdscr.addstr(height-1, 0, statusbarstr)
        stdscr.addstr(height-1, len(statusbarstr), " " * (width - len(statusbarstr) - 1))
        stdscr.attroff(curses.color_pair(3))

        # Turning on attributes for title
        stdscr.attron(curses.color_pair(2))
        stdscr.attron(curses.A_BOLD)

        # Rendering title
        stdscr.addstr(start_y, start_x_title, title)

        # Turning off attributes for title
        stdscr.attroff(curses.color_pair(2))
        stdscr.attroff(curses.A_BOLD)

        # Print rest of text
        for line in asciiArt:
            stdscr.addstr(start_y_art, start_x_art, line)
            start_y_art += 1
        stdscr.addstr(start_y + 1, start_x_subtitle, subtitle)
        stdscr.move(cursor_y, cursor_x)

        # Refresh the screen
        stdscr.refresh()

        # wait for next input
        k = stdscr.getch()

def main():
    with open('config.json') as f:
        config = json.load(f)

    os.makedirs('tmp', exist_ok=True)

    client = Bot()

    threading.Thread(target=lambda: wrapper(draw_menu)).start()
    
    threading.Thread(target=lambda: client.run(config['token'], bot=False)).start()

    #curses.wrapper(textBoxTest)

if __name__ == "__main__":
    main()
