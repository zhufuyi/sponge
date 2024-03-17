## Automated Testing

| Script  | Description  |
| :------ |:-------|
| [**generate_multi_repo_test.sh**](./generate_multi_repo_test.sh)  | Tests the generation of various types of services (suitable for monolith or multi-repo) and their automatic addition of API code, embedding gorm.model. The success of compiling each service indicates that the test has passed. After the test is completed, execute the script [**clean_multi_repo.sh**](./clean_multi_repo.sh) to delete the generated code files. |
| [**generate_mono_repo_test.sh**](./generate_mono_repo_test.sh)   | Tests the generation of various types of services (suitable for mono-repo) and their automatic addition of API code. The success of compiling each service indicates that the test has passed. After the test is completed, execute the script [**clean_mono_repo.sh**](./clean_mono_repo.sh) to delete the generated code files.                                                   |
| [**generate_auto_test.sh**](./generate_auto_test.sh) | Tests the generation of code for 5 types of services (suitable for monolith or multi-repo), runs the services, automatically requests the API interface, and returns data successfully indicates that the test has passed. After the test is completed, execute the script [**clean_auto.sh**](./clean_auto.sh) to delete the generated code files.                    |
