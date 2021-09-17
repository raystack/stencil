from google.protobuf.message import Message
import requests
from google.protobuf.descriptor_pb2 import FileDescriptorSet
from google.protobuf.message import Message
from google.protobuf.message_factory import GetMessages

class Store:
    def __init__(self):
        self.data = {}
    
    def get(self, name) -> Message:
        return self.data.get(name)
    
    def load(self, url):
        result = requests.get(url, stream=True)
        fds = FileDescriptorSet.FromString(result.raw.read())
        messages = GetMessages([file for file in fds.file])
        self.data.update(messages)
