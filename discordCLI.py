import sys,os
import threading
import curses
import json
import aiofiles
import discord
import asyncio

from curses.textpad import Textbox, rectangle

class Bot(discord.Client):
    def __init__(self):
        super(Bot, self).__init__()
    
    async def on_ready(self):
        l = 1
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
        
        # command logic
        if (k == ord('q')):
            statusbarstr = "Are you sure you want to exit? (N/y)"
            draw_status_bar(stdscr, statusbarstr, width, height)
            k = stdscr.getch()
            if(k == ord('y') or k == ord('Y')):
                break
        elif(k == ord('s')):
            draw_server_list(stdscr)
            k = stdscr.getch()

        # Initialization
        stdscr.clear()
        height, width = stdscr.getmaxyx()
        statusbarstr = "Press 'q' to exit | 's' to choose a server | 'p' to look at private messages"

        # Centering calculations
        start_y = int((height // 2) - 2)

        # Render status bar
        draw_status_bar(stdscr, statusbarstr, width, height)

        # Print title and subtitle
        print_title(stdscr, width, start_y)
        print_sub_title(stdscr, width, start_y)

        # Print ascii art
        print_ascii_art(stdscr, width)

        # refresh 
        stdscr.refresh()

        # wait for next input
        k = stdscr.getch()

    # logout of discord after quitting
    logout()
    curses.wrapper(textBoxTest)

def draw_status_bar(stdscr, statusbarstr, width, height):
    stdscr.attron(curses.color_pair(3))
    stdscr.addstr(height-1, 0, statusbarstr)
    stdscr.addstr(height-1, len(statusbarstr), " " * (width - len(statusbarstr) - 1))
    stdscr.attroff(curses.color_pair(3))
    stdscr.refresh()

def print_title(stdscr, width, start_y):
    title = "Discord CLI"[:width-1]
    # centering calculation
    start_x_title = int((width // 2) - (len(title) // 2) - len(title) % 2)
    # setting text attributes
    stdscr.attron(curses.color_pair(2))
    stdscr.attron(curses.A_BOLD)
    # print title
    stdscr.addstr(start_y, start_x_title, title)
    # unset text attributes
    stdscr.attroff(curses.color_pair(2))
    stdscr.attroff(curses.A_BOLD)

def print_sub_title(stdscr, width, start_y):
    subtitle = "Written by Nathan Withers"[:width-1]
    # centering calculation
    start_x_subtitle = int((width // 2) - (len(subtitle) // 2) - len(subtitle) % 2)
    # print subtitle
    stdscr.addstr(start_y + 1, start_x_subtitle, subtitle)

def print_ascii_art(stdscr, width):
    # read in ascii to list
    with open("ascii60Wide.txt") as f:
        content = f.readlines()
        asciiArt = [x.strip() for x in content]
    # centering calculation
    start_x_art = int((width // 2) - (len(asciiArt[1]) // 2) - len(asciiArt[1]) %2)
    start_y_art = 1
    # print ascii art
    for line in asciiArt:
        stdscr.addstr(start_y_art, start_x_art, line)
        start_y_art += 1

def draw_server_list(stdscr):
    y, x = 1, 100
    for server in client.servers:
        stdscr.addstr(y, x, server.name)
        y += 1
        for channel in server.channels:
            stdscr.addstr(y, x, chr(y + 64) + ") " + channel.name)
            y += 1
    stdscr.refresh()

def logout():
    loop = asyncio.get_event_loop()
    loop.run_until_complete(client.logout())
    loop.close()

def textBoxTest(stdscr):
    stdscr.addstr(0, 0, "Enter IM message: (hit Ctrl-G to send)")
    
    width = 30
    height = 5
    
    editwin = curses.newwin(height,width, 2,1)
    rectangle(stdscr, 1,0, 1+height+1, 1+width+1)
    stdscr.refresh()

    box = Textbox(editwin)

    # Let the user edit until Ctrl-G is struck.
    box.edit()

    # Get resulting contents
    message = box.gather()

client = Bot()

def main():
    with open('config.json') as f:
        config = json.load(f)

    #os.makedirs('tmp', exist_ok=True)
    threading.Thread(target=lambda: curses.wrapper(draw_menu)).start()

    loop = asyncio.get_event_loop()
    loop.run_until_complete(client.start(config['token'], bot=False))
    loop.close()
    
if __name__ == "__main__":
    main()
