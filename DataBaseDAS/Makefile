all: database

database: delete insert readJson select supportFiles where main
	g++ -std=c++20 -pthread delete.o insert.o readJson.o select.o supportFiles.o where.o main.o -o databaseresult

main:
	g++ -std=c++20 -c main.cpp

delete:
	g++ -std=c++20 -c ./func/delete.cpp

insert:
	g++ -std=c++20 -c ./func/insert.cpp

readJson:
	g++ -std=c++20 -c ./func/readJson.cpp

select:
	g++ -std=c++20 -c ./func/select.cpp

supportFiles:
	g++ -std=c++20 -c ./func/supportFiles.cpp

where:
	g++ -std=c++20 -c ./func/where.cpp

clean: 
	rm -rf *.o databaseresult

rebuild: clean database
