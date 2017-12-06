import sys
import os
import queue
import threading
import curses
import json
import discord
from discord import ChannelType
import asyncio
import aiofiles

from curses.textpad import Textbox, rectangle

command_queue = queue.Queue()

class Bot(discord.Client):
    def __init__(self):
        super(Bot, self).__init__()
    
    async def poll_queue(self):
        while(command_queue.empty()):
            continue
        if(command_queue.get() == 'logout'):
            await self.logout()
        else:
            await self.poll_queue()
    
    async def on_ready(self):
        #print('Logged in as')
        #print(self.user.name)
        #print(self.user.id)
        #print('------')
        await self.poll_queue()

    async def on_message(self, message):
        if message.author == self.user:
            content = message.content  # type: str

            if content.startswith('$test'):
                await self.send_message(message.channel, 'hmmmm')

client = Bot()

def draw_menu(stdscr):
    k = 0

    # Clear and refresh the screen for a blank canvas
    stdscr.clear()
    stdscr.refresh()

    height, width = stdscr.getmaxyx()
    # Start colors in curses
    curses.start_color()
    curses.init_pair(1, curses.COLOR_CYAN, curses.COLOR_BLACK)
    curses.init_pair(2, curses.COLOR_RED, curses.COLOR_BLACK)
    curses.init_pair(3, curses.COLOR_BLACK, curses.COLOR_WHITE)

    # Loop where k is the last character pressed
    while True: 
        # command logic
        if k == ord('q'):
            statusbarstr = "Are you sure you want to exit? (N/y)"
            draw_status_bar(stdscr, statusbarstr, height, width)
            k = stdscr.getch()
            if k == ord('y') or k == ord('Y'):
                end_curses(stdscr)
                break
        elif k == ord('s'):
            guilds = draw_server_list(stdscr, height, width)
            k = stdscr.getch()
            if chr(k) in guilds:
                draw_splash_screen(stdscr, height, width)
                draw_channel_list(stdscr, guilds[chr(k)], height, width)
                k = stdscr.getch()
        elif k == ord('p'):
            draw_private_channels(stdscr, height, width)
            k = stdscr.getch()

        # redraw splash screen
        draw_splash_screen(stdscr, height, width)

        # wait for next input
        k = stdscr.getch()

    # logout of discord after quitting
    command_queue.put('logout')
    #curses.wrapper(textBoxTest)

def draw_status_bar(stdscr, statusbarstr, height, width):
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
    stdscr.addstr(start_y, start_x_subtitle, subtitle)

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
    return start_y_art

def draw_splash_screen(stdscr, height, width):
    # Initialization
    stdscr.clear()
    statusbarstr = "Press 'q' to exit | 's' to choose a server | 'p' to look at private messages"

    # Render status bar
    draw_status_bar(stdscr, statusbarstr, height, width)

    # Print ascii art
    start_y = print_ascii_art(stdscr, width)

    # Print title and subtitle
    print_title(stdscr, width, start_y + 5)
    print_sub_title(stdscr, width, start_y + 6)
    
    # refresh 
    stdscr.refresh()

def draw_server_list(stdscr, height, width):
    # dictionary of guilds
    guilds = {}
    y, x = 0, 1
    for server in client.servers:
        y += 1
        # draw each guild with a letter for selection
        stdscr.addstr(y, x, chr(y + 64) + ") " + server.name)
        # add guild to dictionary with selector letter as key
        guilds.update({chr(y + 96):server})
    # update status bar
    draw_status_bar(stdscr, "Select a guild using " + chr(65) + "-" + chr(y + 64), height, width)
    stdscr.refresh()
    return guilds

def draw_channel_list(stdscr, guild, height, width):
    # dictionary of channels
    channels = {}
    y, x = 0, 1
    stdscr.addstr(y + 1, x, guild.name)
    for channel in guild.channels:
        if channel.type is ChannelType.text:
            if channel.permissions_for(channel.server.me).read_message_history:
                y += 1
                # draw each channel with a letter for selection
                stdscr.addstr(y + 1,x, chr(y + 64) + ") " + channel.name)
                # add channel to dictionary with selector letter as key
                channels.update({chr(y + 96):channel})
    # update status bar
    draw_status_bar(stdscr, "Select a channel using " + chr(65) + "-" + chr(y + 64), height, width)
    stdscr.refresh()
    return channels

def draw_private_channels(stdscr, height, width):
    # dictionary of channels
    channels = {}
    y, x = 0, 1
    for channel in client.private_channels:
        y += 1
        # draw each channel with a letter for selection
        stdscr.addstr(y,x, chr(y + 64) + ") " + (channel.user.name if channel.name is None else channel.name))
        # add channel to dictionary with selector letter as key
        channels.update({chr(y + 96):channel})
    # update status bar
    draw_status_bar(stdscr, "Select a channel using " + chr(65) + "-" + chr(y + 64), height, width)
    stdscr.refresh()
    return channels

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

def end_curses(stdscr):
    curses.nocbreak()
    stdscr.keypad(False)
    curses.echo()
    curses.endwin()

def main():
    with open('config.json') as f:
        config = json.load(f)

    #os.makedirs('tmp', exist_ok=True)
    threading.Thread(target=lambda: curses.wrapper(draw_menu)).start()
    loop = asyncio.get_event_loop()
    loop.run_until_complete(client.start(config['token'], bot=False))
    loop.close()

if __name__ == '__main__':
    main()
