import sys
import logging

# Print the Python path for debugging purposes
print(sys.path)

try:
    import hcl2
except ImportError as e:
    print(f"Error importing hcl2: {e}")
    print(f"Python version: {sys.version}")
    print(f"Python executable: {sys.executable}")
    sys.exit(1)

def parse_hcl_file(file_path):
    try:
        with open(file_path, 'r') as file:
            content = file.read()
            print(f"File content:\n{content}")  # Debugging statement
            parsed_hcl = hcl2.loads(content)
            print(f"Parsed HCL: {parsed_hcl}")  # Debugging statement
            return parsed_hcl.get('resource', [])
    except FileNotFoundError:
        logging.error(f"Error: File not found: {file_path}")
        sys.exit(1)
    except hcl2.HclParseError as e:
        logging.error(f"Error parsing HCL file: {e}")
        sys.exit(1)
    except Exception as e:
        logging.error(f"Unexpected error parsing HCL file: {e}")
        sys.exit(1)