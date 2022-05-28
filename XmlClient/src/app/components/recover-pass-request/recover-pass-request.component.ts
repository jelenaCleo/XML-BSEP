import { Component, OnInit } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { NgForm } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { UserService } from 'src/app/services/user-service/user.service';

@Component({
  selector: 'app-recover-pass-request',
  templateUrl: './recover-pass-request.component.html',
  styleUrls: ['./recover-pass-request.component.css']
})
export class RecoverPassRequestComponent implements OnInit {

  constructor(
    private authService: UserService,
    private _snackBar: MatSnackBar,
    private _router: Router
    ) { }

  ngOnInit(): void {
  }

  onSubmit(f: NgForm) {

    const loginObserver = {
      next: (x: any) => {
        console.log(x);
        this._router.navigate(['recover']);
        this._snackBar.open("Your code is sent!", "Dismiss");
      },
      error: (err: HttpErrorResponse) => {
       
        this._snackBar.open(err.error.message + "!", 'Dismiss');
      },
    };
    //this.authService.login(f.value).subscribe(loginObserver);
    this._router.navigate(['recover']);
  }

}
