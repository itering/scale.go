import json
import os
import sys
from collections import OrderedDict


def convert_type_registry(name):
    module_path = os.path.dirname(__file__)
    path = os.path.join(module_path, '{}'.format(name))

    with open(os.path.abspath(path), 'r') as fp:
        data = fp.read()
    type_registry = json.loads(data, object_pairs_hook=OrderedDict)

    convert_dict = dict()

    for type_string, type_struct in type_registry.items():

        if type(type_struct) == OrderedDict:
            if '_enum' in type_struct:
                convert_dict[type_string] = {"type": "enum"}
                if type(type_struct["_enum"]) == list:
                    convert_dict[type_string]["value_list"] = try_struct_convert(type_struct["_enum"])
                else:
                    convert_dict[type_string]["type_mapping"] = try_struct_convert(type_struct["_enum"])

            elif '_set' in type_struct:
                set_element = OrderedDict(type_struct["_set"]).keys()
                set_element.pop(0)
                convert_dict[type_string] = {
                    "type": "set",
                    "bit_length": type_struct["_set"]["_bitLength"],
                    "type_mapping": set_element,
                }
                print
            else:
                convert_dict[type_string] = {
                    "type": "struct",
                    "type_mapping": try_struct_convert(type_struct)
                }
        else:
            convert_dict[type_string] = type_struct
    return convert_dict


def try_struct_convert(struct):
    n = []
    if type(struct) == list:
        return struct
    for type_string, type_struct in struct.items():
        if type_struct is None:
            type_struct = "null"
        n.append([type_string, type_struct])
    return n


if __name__ == '__main__':
    covert = convert_type_registry(sys.argv[1])
    print(json.dumps(covert, indent=2))
