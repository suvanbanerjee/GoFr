import os

# Directory containing the files to concatenate
directory = '/Users/yuktha/Downloads/gofr-development/docs'
output_file = 'output.txt'

# Open the output file in write mode
with open(output_file, 'w') as outfile:
    # Iterate over each file in the directory
    for filename in os.listdir(directory):
        # Get the full file path
        file_path = os.path.join(directory, filename)
        
        # Only process files (not directories)
        if os.path.isfile(file_path):
            with open(file_path, 'r') as infile:
                # Write the contents of the current file to the output file
                outfile.write(infile.read())
                outfile.write("\n")  # Add a newline between file contents

print(f"All files in {directory} have been concatenated into {output_file}")
