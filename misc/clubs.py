import json

from jinja2 import Environment, FileSystemLoader

# 1. Create a list of clubs (example with 3 clubs)
with open("./data/clubs.json") as fh:
    clubs = json.load(fh)

# 2. Set up Jinja2 environment, assuming your template is in the same folder
env = Environment(loader=FileSystemLoader('./misc'))

# 3. Load the template
template = env.get_template('templates/clubs.html')  # filename of your template

# 4. Render the template with the club data
rendered_html = template.render(clubs=sorted(clubs, key=lambda c: c["name"]))

# 5. Optionally, write the rendered HTML to a file
with open('./misc/clubs.html', 'w', encoding='utf-8') as f:
    f.write(rendered_html)

print("HTML file generated: clubs.html")