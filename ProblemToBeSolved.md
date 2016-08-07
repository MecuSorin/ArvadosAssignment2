## ASSIGNMENT 2

# Given
- a block of data (up to 64 MiB) 
- an ordered list of servers

# Action 
- send asynchronous data to 2 servers using POST

# Remarks
Send the data to the first two servers using two asynchronous POST requests. 
If one fails, 
    let the other continue while initiating a third request on the
    third server. 
Continue until two different servers have reported success.

# Expectations
To test your client, write a stub server that accepts data, 
discards the data,
 and reports HTTP 200, say, 30% of the time, and only after a variable delay.
 
# Build
Use 'go get' and 'go build' to add the required libraries and create the executable
 
# Usage
ArvadosAssignment2 --help will provide the actualized flags

'ArvadosAssignment2 -stubPercent 30 -stub true'
will start as a stub server and will add in servers.json current instance

'ArvadosAsignment2'
will try to send the data to the instances from servers.json list