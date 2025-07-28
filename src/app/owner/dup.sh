#!/bin/bash
#to make it excutable: chmod +x dup.sh
# Specify the source directory
source_dir="user"
source_Cap="User"

# Specify the list of destination directories
destination_dirs=( "inviteCode"  )
#destination_dirs=("user" "admin"   )

# Iterate through the destination directories
for destination_dir in "${destination_dirs[@]}"
do
    # Create the destination directory if it doesn't exist
    mkdir -p "$destination_dir"

    # Copy all files from the source directory to the destination directory
    cp -r "$source_dir"/* "$destination_dir"/

    # Rename files and variables within the destination directory
    cd "$destination_dir" || exit
    for file in *
    do
        # Rename files
        mv "$file" "$(echo "$file" | sed "s/$source_dir/$destination_dir/")"

        # Rename variables (assuming case-sensitive replacements)
        sed -i '' "s/$source_dir/$destination_dir/g" *
        # Capitalize first letter
        first_char=$(echo "${destination_dir:0:1}" | tr '[:lower:]' '[:upper:]')
        rest_chars="${destination_dir:1}"
        new_word="$first_char$rest_chars"
        sed -i '' "s/$source_Cap/$new_word/g" *
#        sed -i '' "s/Domain/${destination_dir^}/g" *  # Capitalize first letter
    done

    cd ..  # Go back to the parent directory
done
echo "Variable replacements complete!"