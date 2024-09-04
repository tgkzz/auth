proto-get:
	protoc -I proto proto/auth/auth.proto --go_out=./gen/go/ --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative

#-I proto: Опция -I или --proto_path указывает путь к корневой директории с .proto файлами. Это нужно для того, чтобы компилятор смог найти импорты, если они есть. В нашем случае это директория proto.
#proto/sso/sso.proto: путь к конкретному .proto файлу, который мы компилируем.
#--go_out=./gen/go/: Опция --go_out указывает, куда записывать сгенерированный Go-код. В нашем случае — ./gen/go/.
#--go_opt=paths=source_relative: дополнительная опция — указывает, как создавать имена пакетов. paths=source_relative означает, что выходные файлы будут иметь тот же пакет, что и исходные .proto файлы.
#--go-grpc_out=./gen/go/: куда записывать сгенерированный Go gRPC-код. Как и в предыдущем случае, выходные файлы будут помещены в директорию ./gen/go/.
#--go-grpc_opt=paths=source_relative: Это аналогичная опция для генерации Go gRPC-кода, указывающая, как создавать имена пакетов для gRPC.