import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { IComment } from 'src/app/interfaces/comment';
import { IInterview } from 'src/app/interfaces/interview';
import { IJobOffer } from 'src/app/interfaces/job-offer';
import { IJobOfferRequest } from 'src/app/interfaces/job-offer-request';
import { IJobOfferResponse } from 'src/app/interfaces/job-offer-response';
import { NewCompanyRequestDto } from 'src/app/interfaces/new-company-request-dto';
import { ISalaryComment } from 'src/app/interfaces/salary-comment';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class CompanyService {

  CreateJobOffer(jobOffer: IJobOfferRequest): Observable<any> {
    return this.http.post(`${this.apiServerUrl}/company/createOffer`, jobOffer);
  }
  UpdateInfo(company: any): Observable<any> {
    return this.http.put(`${this.apiServerUrl}/company/edit/` + company.companyId
      , company);
  }

  private apiServerUrl = environment.apiBaseUrl;
  RegisterCompany(company: NewCompanyRequestDto): Observable<any> {
    return this.http.post(`${this.apiServerUrl}/company/new`, company, {
      responseType: 'text',
    });

  }

  constructor(private http: HttpClient) { }

  getAlCompaniesForUser(): Observable<any> {
    return this.http.get(`${this.apiServerUrl}/company/getAllForUser`);

  }

  getAllPendingCompanies(): Observable<any> {
    return this.http.get(`${this.apiServerUrl}/company/pending`);

  }

  approveRequest(id: number): Observable<any> {
    return this.http.get(`${this.apiServerUrl}/company/approve/` + id);
  }

  rejectRequest(id: number): Observable<any> {
    return this.http.get(`${this.apiServerUrl}/company/reject/` + id);
  }
  getAllUsersCompanies(username: string): Observable<any> {
    return this.http.get(`${this.apiServerUrl}/company/users-company/` + username);
  }

  getById(id: any): Observable<any> {
    return this.http.get(`${this.apiServerUrl}/company/` + id);
  }
  getOffersForCompany(id: string): Observable<IJobOffer[]> {
    return this.http.get<IJobOffer[]>(`${this.apiServerUrl}/offer/all/` + id);
  }

  getCommentsForCompany(id: string): Observable<IComment[]> {
    return this.http.get<IComment[]>(`${this.apiServerUrl}/company/` + id + `/comments`);
  }
  getInterviewsForCompany(id: string): Observable<IInterview[]> {
    return this.http.get<IInterview[]>(`${this.apiServerUrl}/company/` + id + `/interviews`);
  }
  getSalaryCommentsForCompany(id: string): Observable<ISalaryComment[]> {
    return this.http.get<ISalaryComment[]>(`${this.apiServerUrl}/company/` + id + `/salaryComments`);
  }


}
