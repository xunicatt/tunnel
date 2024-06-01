all: tunnel clean

project=tunnel
version=1.0.0
install_dir=${HOME}/.$(project)
compiler=go

tunnel: main.go
	$(compiler) build -o $(project) .

run:
	./$(project)

clean:
	rm -f ./$(project)

install: tunnel
	mkdir -p $(install_dir)
	mkdir -p $(install_dir)/bin $(install_dir)/pkg $(install_dir)/cache $(install_dir)/lib
	cp ./$(project) $(install_dir)/bin
	cp ./contents/init.json $(install_dir)
	-cp -n ./contents/list.json $(install_dir)
	-rm ./$(project)
	@echo
	@echo -e "\033[1;32mAdd '\$$HOME/.tunnel/bin' to your \$$PATH enviorment varibale using:" 
	@echo -e "		export PATH=\$$PATH:\$$HOME/.tunnel/bin\033[0m"
	@echo

uninstall:
	sudo rm -f $(install_dir)/$(project)
