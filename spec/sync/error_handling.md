# Error Handling

We need to be able to handle cases where the execution client responds with an "error'. I use 
the word error in quotes since these are slightly expected behaviour. Notably Status_SYNCING
Status_ACCEPTED and sometimes Status_INVALID. We need retry logic and proper state management to ensure that the beacon change can gracefully handle these edge cases.

