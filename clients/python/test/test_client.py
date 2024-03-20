from unittest.mock import patch

from google.protobuf.reflection import GeneratedProtocolMessageType

from src.raystack.stencil import Client
from src.raystack.store import Store

URL = 'http://stencil.test/proto-descriptors/test/latest'


def get_file_desc():
    with open('data/one.desc', 'rb') as myfile:
        desc = myfile.read()
    return desc


def test_store(protoc_setup):
    file_desc = get_file_desc()
    store = Store()
    store.load(data=file_desc)
    assert 'test.One' in store.data
    assert isinstance(store.get('test.One'), store.get('test.One').__class__)


@patch('raystack.store.Store._load_from_url')
def test_client(test_desc_from_url, protoc_setup):
    file_desc = get_file_desc()
    test_desc_from_url.return_value = file_desc

    client = Client(URL)

    assert isinstance(client.get_descriptor('test.One'), GeneratedProtocolMessageType)
