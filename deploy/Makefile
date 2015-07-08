all:
	cat cloudformation.yml | ruby -e 'require "yaml"; require "json"; print JSON.pretty_generate(YAML.load(STDIN),indent:"  ")' > cloudformation.json