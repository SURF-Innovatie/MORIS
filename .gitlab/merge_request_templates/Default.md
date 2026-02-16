## Related Issues

<!-- Link the GitLab issues related to this MR using closing patterns -->
<!-- Example: Closes #42, Relates to #15 -->

## Description

<!-- Provide a detailed explanation of the changes you have made. Include the reasons behind these changes and any relevant context. -->

## Side Effects

<!-- Describe any side effects of your changes. This might include changes to the behavior of existing features, or changes to the performance of the application. Write "None" if not applicable. -->

## Database

<!-- If your changes affect the database, describe the changes you have made. This might include new migrations, changes to existing migrations, or changes to the database schema. Write "N/A" if not applicable. -->

## Additional Information

<!-- Any additional information that reviewers should be aware of. -->

## Definition of done

Let's make sure we don't forget anything!
Please review the following list and check the boxes that apply to this MR.

If something is not applicable, please check the box anyway.
If something is not possible, please leave a comment explaining why.

- [ ] Code satisfies the requirement as specified in the linked issue
- [ ] The not-so-happy flow and errors are handled gracefully
- [ ] Project builds without errors or warnings
- [ ] Any decisions or changes to the requirements have been added to the issue description
- [ ] Access control is in compliance with the Principle of Least Privilege (PoLP)
- [ ] Code is well-commented and (JS/Go)Docs have been updated where needed
- [ ] Added ENV variables have been added to the `.env.example`
- [ ] Manual changes that need to be performed to the database or ENV have been communicated to the release manager
- [ ] Make sure existing production data is migrated/updated if necessary
- [ ] (Manually) Tested if changes might have caused regression bugs
- [ ] (Manually) Tested if changes might have unintended results for various roles
- [ ] (Manually) Tested UI on desktop, tablet and mobile screen sizes
  - UI is compatible for screens with a minimum of 340px width
- [ ] (Manually) Tested if changes might have unintended results for various browsers
  - Chrome (or Chromium-based browsers)
  - Safari
  - Firefox
