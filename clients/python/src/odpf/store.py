import requests
from google.protobuf.descriptor_pb2 import FileDescriptorSet
from google.protobuf.descriptor import Descriptor
from google.protobuf.message_factory import GetMessages

class Store:
    def __init__(self):
        self.data = {}
    
    def get(self, name) -> Descriptor:
        return self.data.get(name)
    
    def load(self, url):
        result = requests.get(url)
        fds = FileDescriptorSet.FromString(result.text)
        messages = GetMessages([file for file in fds.file])
        self.data.update(messages)
        

        
        