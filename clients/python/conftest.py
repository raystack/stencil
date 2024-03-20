from subprocess import run

import pytest
import os


@pytest.fixture(scope="session")
def protoc_setup():
    current_dir = os.path.dirname(os.path.realpath(__file__))
    output_file = os.path.join(current_dir, 'test/data/one.desc')

    input_dir = os.path.join(current_dir, 'test/data')
    run(['protoc', f'--descriptor_set_out={output_file}', '--include_imports', f'--proto_path={input_dir}',
         'one.proto'], cwd=current_dir)
