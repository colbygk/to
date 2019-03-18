
## Text file to track todo/notes

Parse a text/notes file
  The text file is a format that is both human and machine readable
  todoML so to speak
  ASCII originating file
    can output to richer formats
  Each line may or may not start with white-space
    The amount of white space may or may not determine the relationship of the current line to the previous line
  Each line must have a meaningful indicator character before the remainder of the line
  White space must exist between the indicator character and the rest of the text on the line
  '-' start of a tree/regular text
  '/' this line has been "completed"
  'x' this line will not be completed/is obsolete
  '\' this line is a continuation of the text from the previous line, usually this is done for formatting reasons, so, attempt to honor the basic text formatting as it applies to the previous line
  'r' this line has been re-edited and is logically the same as the next line, but contains the old text
  '!' same as '-' for interpretation but indicates that this is an immediate item to work on
  '/' same as '-' for interpretation but indicates that this is "done/completed"
  's' same as '-' for interpretation but indicates that this is "asleep" 
  'b' same as '-' for interpretation but indicates that this is blocked on something else
      \ maybe this should only ever be a leaf - so that the block can be specific to some other item?
      \ really all of these are "same as '-' for interpretation ..."

  Self-documenting example:
```
- Some project/info trying to work on/todo
  - Something important about it
    \ continuation of the previous thing but wanted to break it
    \ across lines in a pleasant to read way
  x thought this part was important, but it's not, so deleted it
  / 
  r I thought I meant one thing
  - But I really meant another
```

  Real-world example:

```
- homepage row generation failing even with extended time
  - wth
  - Added more logging on virt1 only
    - Did find it will block when backups are being taken because of
      \ "waiting for table to flush" and copying to tmp tables on disk
    - Looks like this is focused on PROJECT_ROWS:filtering:community_most_remixed_projects
- sentry.io
  - scratch-clouddata
    - Rework sentry sending here first
  / Updates to scratchr2 to turn off default event sending
  r consolidating from scratch-foundation to scratchr2
  / consolidating from scratch-team to scratch-foundation
  r Update API keys to point at older setup
  / Update API keys to point at newer setup
    x scratch-www/scratch-gui
    x scratch-assets
    x scratch-projects
    x scratch-api
```

  * `- homepage ...` is a new 'tree' or 'grouping' of related text/notes.
  * `  - wth` with the extra white space on this line, `wth` is now related to previous line
  * `  - Added more ...` indicates another added text/note to the `- homepage ...` tree, at the same level as `  - wth`
  * `    - Did find it will ...` indicates a leaf of `  - Added more ...`
  * `      \ "waiting for ...` adds to the previous line while making it more presentable


# Program that can make use of these notes
  Written in go as `cgk.sh/to`
  Facilitate absurd commands like `to day` and `to morrow` and `to ruminate` a synonym for `to chew`

  - `to day`
    - Finds most recent `tracking.to` file (format described above)
    - Copies it to a new "day", by default it is $HOME/notes/YYYY/mon/DD/tracking.to
    - Prints out the list of "top level" items in the file
      - Should they have numbers? That way they can be expanded with say 'to do ?N'
      - 'to do' synonym for 'to day' with 'N' allows for diving into tree
    - Also constructs a DB file that encodes more information about each item.
      - creates an sha1 hash
        - hash is only on the text itself, but probably will need some other uniqueifying information such as its root's hash and leaf hashes?
          \ this could allow for time stamping in the DB and the ability to track an item
          \ Might be able to parse date/times out of a file to do checking for stale/asleep items
    

