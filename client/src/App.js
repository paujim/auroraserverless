import React from 'react';
import Container from '@material-ui/core/Container';
import Box from '@material-ui/core/Box';
import Paper from '@material-ui/core/Paper';
import ProfilesTable2 from './components/ProfilesTable2'
import AppMenu from './components/AppMenu'
import Snackbar from '@material-ui/core/Snackbar';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';


export default function App() {
  const [error, setError] = React.useState(null);

  const handleClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    setError(null);
  };

  return (
    <Box>
      <AppMenu />
      <Container maxWidth="md">
        <Paper elevation={3} >
          <ProfilesTable2 showError={setError} />
        </Paper>
        <Snackbar
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'left',
          }}
          open={error != null}
          autoHideDuration={6000}
          onClose={handleClose}
          message={error}
          action={
            <React.Fragment>
              <IconButton size="small" aria-label="close" color="inherit" onClick={handleClose}>
                <CloseIcon fontSize="small" />
              </IconButton>
            </React.Fragment>
          }
        />
      </Container>
    </Box>
  );
}