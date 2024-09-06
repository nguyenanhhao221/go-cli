#!/bin/bash

# Calculate the initial hash of the file using md5
FHASH=$(md5 -q "$1")

while true; do
    # Calculate the current hash of the file
    NHASH=$(md5 -q "$1")
    
    # Check if the hash has changed
    if [ "$NHASH" != "$FHASH" ]; then
        # Run the command if the file has changed
        ./mdp -file "$1"
        
        # Update the hash
        FHASH=$NHASH
    fi
    
    # Wait for 5 seconds before checking again
    sleep 5
done
