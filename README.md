# TldrBot #
TL;DR bot for Telegram

This is my attempt to write Telegram bot in Go. TL;DR stands for "Too long; didn't read" if you're wondering.

Typical workflow:

1. This bot is invited to some chat
2. Someone sends a link to article
3. Bot intercepts this message (don't worry, it skips all other messages)
4. Bot retrieves the page that the link points too
5. Bot tries to find main content area of that page
6. Bot extracts most meaningful phrases from that article
7. Bot then proceeds to publish this chewed content as a message to the chat where link appeared

# Libraries used #

* [TLDR text summarizer](https://github.com/JesusIslam/tldr) (MIT License)
* [Html Content / Article Extractor](github.com/advancedlogic/GoOse) (Apache License 2.0)

# License #

Copyright (C) 2016  Oleg Chernovskiy

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
