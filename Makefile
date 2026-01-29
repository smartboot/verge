reset_git:
	git checkout --orphan new_branch
	git add -A
	git reset -- .idea config vendor res .driverbox_serial_no verge-export _output main
	git commit -m "Initial commit"
	git branch -D master
	git branch -m master
	git push -f origin master