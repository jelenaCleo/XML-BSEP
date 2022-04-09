package com.example.PKI.service.cert.impl;

import com.example.PKI.dto.CertificateDto;
import com.example.PKI.model.CertificateType;
import com.example.PKI.model.Subject;
import com.example.PKI.model.User;
import com.example.PKI.repository.CertificateRepository;
import com.example.PKI.repository.UserRepository;
import com.example.PKI.service.KeyService;
import com.example.PKI.service.KeyStoreService;
import com.example.PKI.service.cert.CertificateService;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x500.X500NameBuilder;
import org.bouncycastle.asn1.x500.style.BCStyle;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509CertificateHolder;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.OperatorCreationException;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.math.BigInteger;
import java.security.*;
import java.security.cert.Certificate;
import java.security.cert.CertificateEncodingException;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.time.LocalDate;
import java.util.Date;

@Service
public class CertificateServiceImpl implements CertificateService {

    @Autowired
    private KeyService keyService;
    @Autowired
    private KeyStoreService keyStoreService;
    @Autowired
    private CertificateRepository certificateRepository;
    @Autowired
    private UserRepository userRepository;
    private BigInteger serialNumber;

    @Override
    public X509Certificate generateCertificate(CertificateDto certificateDto, Subject generatedSubjectData) throws Exception {
        try {

            if (LocalDate.parse(certificateDto.getEndDate()).compareTo(LocalDate.parse(certificateDto.getStartDate())) < 0) {
                throw new Exception("END_DATE_BEFORE_START");
            }

//            User issuer = userRepository.findById(certificateDto.getIssuerId()).get();
//            if (issuer == null) {
//                throw new Exception("ISSUER_NOT_FOUND");
//            }

            String issuerAlias = certificateDto.getIssuerSerialNumber();

            if (certificateRepository.findTypeBySerialNumber(certificateDto.getIssuerSerialNumber()).equals(CertificateType.ROOT)) {

            }

            SimpleDateFormat sdf = new SimpleDateFormat("yyyy-MM-dd");
            Date startDate = sdf.parse(certificateDto.getStartDate());
            Date endDate = sdf.parse(certificateDto.getEndDate());

            serialNumber = new BigInteger(keyService.getSerialNumber().toString());
            //ovo istraziti : on koristi SHA256withECDSA
            JcaContentSignerBuilder builder = new JcaContentSignerBuilder("SHA256WithRSAEncryption");
            Security.addProvider(new BouncyCastleProvider());
            builder.setProvider("BC");
            KeyStore keyStore = keyStoreService.getKeyStore(keyService.getKeyStorePath(certificateDto.getType()), keyService.getKeyStorePass());

            X500Name issuer = null;
            PrivateKey privateKey = null;
            if (certificateDto.getType().equals("ROOT")) {
                issuer = generatedSubjectData.getX500Name();
                privateKey = generatedSubjectData.getKeyPair().getPrivate();
            } else {
                User issuerData = userRepository.findById(certificateDto.getIssuerId()).get();
                X509Certificate issuerCert = (X509Certificate) keyStore.getCertificate(issuerData.getEmail() + issuerData.getId());
                issuer = new JcaX509CertificateHolder(issuerCert).getSubject();
                privateKey = (PrivateKey) keyStore.getKey(certificateDto.getIssuerSerialNumber(), "key".toCharArray());
            }

            ContentSigner contentSigner = builder.build(privateKey);

            X509v3CertificateBuilder certGen = new JcaX509v3CertificateBuilder(issuer,
                    serialNumber,
                    startDate,
                    endDate,
                    generatedSubjectData.getX500Name(),
                    generatedSubjectData.getKeyPair().getPublic());

            X509CertificateHolder certHolder = certGen.build(contentSigner);

            JcaX509CertificateConverter certConverter = new JcaX509CertificateConverter();
            BouncyCastleProvider bcp = new BouncyCastleProvider();
            certConverter = certConverter.setProvider(bcp);

            //Konvertuje objekat u sertifikat
            return certConverter.getCertificate(certHolder);

        } catch (CertificateEncodingException e) {
            e.printStackTrace();
        } catch (IllegalArgumentException e) {
            e.printStackTrace();
        } catch (IllegalStateException e) {
            e.printStackTrace();
        } catch (OperatorCreationException e) {
            e.printStackTrace();
        } catch (CertificateException e) {
            e.printStackTrace();
        } catch (UnrecoverableKeyException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        } catch (KeyStoreException e) {
            e.printStackTrace();
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
        } catch (ParseException e) {
            e.printStackTrace();
        }

        return null;
    }

    @Override
    public com.example.PKI.model.Certificate saveCertificateDB(CertificateDto subjectDto, Integer subjectId) {
        CertificateType cerType;
        if (subjectDto.getType().equals("ROOT")) {
            cerType = CertificateType.ROOT;
        } else if (subjectDto.getType().equals("INTERMEDIATE")) {
            cerType = CertificateType.INTERMEDIATE;
        } else {
            cerType = CertificateType.CLIENT;
        }

        User subject = userRepository.findById(subjectId).get();

        if (subjectDto.getStartDate().compareTo(subjectDto.getEndDate()) < 0) {
            com.example.PKI.model.Certificate certificate = new com.example.PKI.model.Certificate();
            certificate.setSerialNumber(serialNumber.toString());
            certificate.setType(cerType);
            certificate.setValid(true);
            certificate.setSubjectCommonName(subject.getCommonName() + " " + subject.getOrganization());
            certificate.setValidFrom(subjectDto.getStartDate());
            certificate.setValidTo(subjectDto.getEndDate());
            return certificateRepository.save(certificate);
        }
        return null;
    }

    @Override
    public void createCertificate(CertificateDto certificateDto, Subject generatedSubjectData) throws Exception {

        User user = userRepository.findById(certificateDto.getSubjectId()).get();

        X509Certificate certificate = generateCertificate(certificateDto, generatedSubjectData);
        String alias = certificate.getSerialNumber().toString();
        KeyStore keyStore = keyStoreService.getKeyStore(keyService.getKeyStorePath(certificateDto.getType()), keyService.getKeyStorePass());
        keyStoreService.store(keyService.getKeyStorePass(), keyService.getKeyPass(), new Certificate[]{certificate}, generatedSubjectData.getKeyPair().getPrivate(), alias, keyService.getKeyStorePath(certificateDto.getType()));

    }

    @Override
    public X509Certificate getCertificateByAlias(String alias, KeyStore keystore) throws KeyStoreException {
        X509Certificate certificate = (X509Certificate) keystore.getCertificate(alias);
        return certificate;
    }

    @Override
    public Subject generateSubjectData(Integer subjectId) {

        User subject = userRepository.findById(subjectId).get();

        KeyPair keyPairSubject = keyService.generateKeyPair();

        X500NameBuilder builder = new X500NameBuilder(BCStyle.INSTANCE);
        builder.addRDN(BCStyle.CN, subject.getCommonName());
        builder.addRDN(BCStyle.O, subject.getOrganization());
        builder.addRDN(BCStyle.OU, subject.getOrganizationUnit());
        builder.addRDN(BCStyle.C, subject.getCountry());
        builder.addRDN(BCStyle.E, subject.getEmail());

        return new Subject(builder.build(), keyPairSubject);
    }


   /* private boolean validate(String alias) throws CertificateException, NoSuchAlgorithmException, KeyStoreException, IOException {
        KeyStore keyStore=keyStoreService.getKeyStore(keyService.getKeyStorePath(),keyService.getKeyStorePass());
        X509Certificate certificate= (X509Certificate) keyStore.getCertificate(alias);
        Certificate[] chain=keyStore.getCertificateChain(alias);
        return false;
    }*/
}
