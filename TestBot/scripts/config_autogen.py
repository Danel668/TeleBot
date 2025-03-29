import json

TEMPLATE = '''/* This file is autogen, don`t edit! */

package sources
{structs}'''

STRUCT_TEMPLATE = '''
type {type_name} struct {{
{fields}
}}
'''

FIELD_TEMPLATE = '    {field_name} {field_type} `json:"{json_name}"`'

def to_go_name(json_name):
    parts = json_name.split('_')
    return ''.join(part.capitalize() for part in parts)


def determine_field_type(v, k):
    if isinstance(v, dict):
        return to_go_name(k)
    elif isinstance(v, str):
        return 'string'
    elif isinstance(v, bool):
        return 'bool'
    elif isinstance(v, int):
        return 'int'
    elif isinstance(v, float):
        return 'float64'
    elif isinstance(v, list):
        return '[]interface{}'
    else:
        return 'interface{}'


def generate_structs(data):
    lines = []
    for key, value in data.items():
        if isinstance(value, dict):
            struct_name = to_go_name(key)
            fields = '\n'.join([FIELD_TEMPLATE.format(
                field_name=to_go_name(k),
                field_type=determine_field_type(v, k),
                json_name=k
            ) for k, v in value.items()])
            struct_def = STRUCT_TEMPLATE.format(type_name=struct_name, fields=fields)
            lines.append(struct_def)
            lines.extend(generate_structs(value))
    return lines


def create_go_file(structs_code):
    with open('src/models/sources/config_types.go', 'w') as f:
        f.write(TEMPLATE.format(structs=''.join(structs_code)))


def main():
    with open('config/production.json', 'r') as f:
        data = {'config': json.load(f)}
    
    struct_definitions = generate_structs(data)

    create_go_file(struct_definitions)

    print("Файл src/models/sources/config_types.go был успешно перезаписан!")


if __name__ == '__main__':
    main()
