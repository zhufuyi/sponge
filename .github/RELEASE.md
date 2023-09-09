## Change log

- Add modify duplicate error codes command, avoid manual modification of duplicate error codes.
- Change `ID` field go type to `uint64` in sql based generated model code to avoid ID type inconsistency.
- Add `make update-config` command.
- Add view error codes list api interface `/codes`.
- Add view service configuration list api interface `/config`.
- Fix Mac M1 install sponge failed ([#8](https://github.com/zhufuyi/sponge/issues/8))
