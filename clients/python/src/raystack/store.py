import requests
from google.protobuf.descriptor_pb2 import FileDescriptorSet
from google.protobuf.message import Message
from google.protobuf.message_factory import GetMessages


class Store:
    def __init__(self):
        self.data = {}

    def get(self, name) -> Message:
        return self.data.get(name)

    def _load_from_url(self, url):
        result = requests.get(url, stream=True)
        return result.raw.read()

    def load(self, url: str = None, data: bytes = None):
        if url:
            data = self._load_from_url(url)
        fds = FileDescriptorSet.FromString(data)
        messages = GetMessages([file for file in fds.file])
        self.data.update(messages)
