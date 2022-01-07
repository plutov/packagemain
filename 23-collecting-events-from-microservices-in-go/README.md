## packagemain #23: Collecting Events from Microservices in Go

- introduction into a real use case, collecting logs centrally from multple microservices
- possible options, collecting stdout, but we're more interested in collecting custom metrics from our services
- idea: create a go grpc service whose responsibbility is to receive events from different services, and be a middleware for storing the events eventually
- pros: defined format, so all services follow the same format, but with possibility to extend easily, single place to maage dependecy with events storage
- how to implement: grpc service, go module for go client, evet structure, storage interface
