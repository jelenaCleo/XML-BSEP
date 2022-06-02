import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { JobOfferComponent } from 'src/app/components/job-offer/job-offer.component';
import { LeaveCommentComponent } from 'src/app/components/leave-comment/leave-comment.component';
import { LeaveInterviewCommentComponent } from 'src/app/components/leave-interview-comment/leave-interview-comment.component';
import { LeaveSallaryCommentComponent } from 'src/app/components/leave-sallary-comment/leave-sallary-comment.component';
import { IComment } from 'src/app/interfaces/comment';
import { CompanyResponseDto } from 'src/app/interfaces/company-response-dto';
import { IInterview } from 'src/app/interfaces/interview';
import { IsUsersCompanyDto } from 'src/app/interfaces/is-users-company-dto';
import { IJobOffer } from 'src/app/interfaces/job-offer';
import { ISalaryComment } from 'src/app/interfaces/salary-comment';
import { CompanyService } from 'src/app/services/company-service/company.service';


@Component({
  selector: 'app-company-profile',
  templateUrl: './company-profile.component.html',
  styleUrls: ['./company-profile.component.css']
})
export class CompanyProfileComponent implements OnInit {
  description!: string;
  editable!: boolean;
  newDescription!: string
  jobOffers!: IJobOffer[]
  comments!: IComment[]
  interviews!: IInterview[]
  salaryComments!: ISalaryComment[]

  company!: CompanyResponseDto;

  cid!: string;

  isUsersCompany! : string;
  role! : string | null;

  constructor(
    private router: Router,
    private _snackBar: MatSnackBar,
    private companyService: CompanyService,
    public matDialog: MatDialog
  ) {

    this.company = {} as CompanyResponseDto;


    this.editable = false;
  
  }

  ngOnInit(): void {
    console.log(this.router.url);
    this.cid = this.router.url.substring(9);
    this.role = localStorage.getItem('role');

    this.companyService.isUsersCompany(this.cid).subscribe(
      res => {
        this.isUsersCompany=res.message;
        console.log(this.isUsersCompany);
      }
    );

    this.companyService.getById(this.cid).subscribe(
      res => {
        this.company = res;
      }
    );

    this.getOffersForCompany()

    this.companyService.getCommentsForCompany(this.cid).subscribe({
      next: (result) => {
        this.comments = result;
      },
      error: (data) => {
        if (data.error && typeof data.error === 'string')
          console.log('desila se greska2');
      },
    });

    this.companyService.getInterviewsForCompany(this.cid).subscribe({
      next: (result) => {
        this.interviews = result;
      },
      error: (data) => {
        if (data.error && typeof data.error === 'string')
          console.log('desila se greska3');
      },
    });

    this.companyService.getSalaryCommentsForCompany(this.cid).subscribe({
      next: (result) => {
        this.salaryComments = result;
      },
      error: (data) => {
        if (data.error && typeof data.error === 'string')
          console.log('desila se greska4');
      },
    });
  }

  getOffersForCompany() {
    this.companyService.getOffersForCompany(this.cid).subscribe({
      next: (result) => {
        this.jobOffers = result;
      },
      error: (data) => {
        if (data.error && typeof data.error === 'string')
          console.log('desila se greska1');
      },
    });
  }

  enableEdit() {
    this.editable = true;
  }
  updateInfo() {
    this.editable = false;

    this.companyService.UpdateInfo(this.company).subscribe({
      next: () => {
        this._snackBar.open(
          'Your request has been successfully submitted.',
          'Dismiss'
        );

      },
      error: (err: HttpErrorResponse) => {
        this._snackBar.open(err.error.message + "!", 'Dismiss');
      },
      complete: () => console.info('complete')
    });

  }
  openModal() {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose = false;
    dialogConfig.id = 'modal-component';
    dialogConfig.height = 'fit-content';
    dialogConfig.width = '500px';
    const dialogRef = this.matDialog.open(JobOfferComponent, dialogConfig);
    dialogRef.afterClosed().subscribe({
      next: (res) => {
        console.log(res);
        this.getOffersForCompany()
      }
    })
  }

  openLeaveComment() {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose = false;
    dialogConfig.id = 'modal-component';
    dialogConfig.height = 'fit-content';
    dialogConfig.width = '500px';
    this.matDialog.open(LeaveCommentComponent, dialogConfig); //TODO: OVDJE JOB OFFER

  }

  openLeaveInterviewComment() {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose = false;
    dialogConfig.id = 'modal-component';
    dialogConfig.height = 'fit-content';
    dialogConfig.width = '500px';
    this.matDialog.open(LeaveInterviewCommentComponent, dialogConfig); //TODO: OVDJE JOB OFFER

  }

  openLeaveSallaryComment() {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose = false;
    dialogConfig.id = 'modal-component';
    dialogConfig.height = 'fit-content';
    dialogConfig.width = '500px';
    this.matDialog.open(LeaveSallaryCommentComponent, dialogConfig); //TODO: OVDJE JOB OFFER

  }


  
}
