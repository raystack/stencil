class MultiUrlClient:
    def __init__(self, urls:list) -> None:
        pass
    
    


class Client(MultiUrlClient):
    def __init__(self, url: str) -> None:
        super().__init__([url])