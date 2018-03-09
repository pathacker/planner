Planner reads the named files, default $HOME/lib/calendar,
and writes to standard output, in calendar order, any lines
containing matching dates for today and tomomrrow.
The '-n days' flag changes the number of days compared.
No special processing is done for weekends.

Recognized date formats are "4/26", "Apr 26", "26 April".
Only the first three runes of the month name are matched.

All comparisions are case insensitive.
