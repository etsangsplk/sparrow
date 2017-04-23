/*
 * Copyright (c) 2016 Keydap Software.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * See LICENSE file for details.
 */
package com.keydap.sparrow;

import java.lang.reflect.Field;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.List;
import java.util.TimeZone;

import org.apache.http.HttpStatus;
import org.junit.BeforeClass;

import com.keydap.sparrow.auth.SparrowAuthenticator;
import com.keydap.sparrow.scim.Device;
import com.keydap.sparrow.scim.Group;
import com.keydap.sparrow.scim.User;
import com.keydap.sparrow.scim.User.Email;
import com.keydap.sparrow.scim.User.Name;

import static org.apache.commons.lang.RandomStringUtils.randomAlphabetic;
import static org.junit.Assert.*;

/**
 *
 * @author Kiran Ayyagari (kayyagari@keydap.com)
 */
public abstract class TestBase {
    
    protected static String baseApiUrl = "http://localhost:7090/v2";
    
    protected static SparrowClient client;
    
    /** the anonymous client */
    static SparrowClient unAuthClient;
    
    static DateFormat utcDf = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'");//RFC3339

    static SparrowAuthenticator authenticator;

    @BeforeClass
    public static void createClient() throws Exception {
        authenticator = new SparrowAuthenticator("admin", "example.COM", "secret");

        client = new SparrowClient(baseApiUrl, authenticator);
        client.register(User.class, Group.class, Device.class);
        utcDf.setTimeZone(TimeZone.getTimeZone("UTC"));

        client.authenticate();
        assertNotNull(authenticator.getToken());
        System.out.println(authenticator.getToken());
        
        unAuthClient = new SparrowClient(baseApiUrl);
        unAuthClient.register(User.class, Group.class, Device.class);
    }
    
    public static <T> void deleteAll(Class<T> resClass) {
        SearchResponse<T> resp = client.searchResource("id pr", resClass, "id", "username");
        List<T> existing = resp.getResources();
        if(existing != null) {
            try {
                Field id = resClass.getDeclaredField("id");
                id.setAccessible(true);
                
                for(T u : existing) {
                    Response<Boolean> delResp = client.deleteResource(id.get(u).toString(), resClass);
                    if(delResp.getHttpCode() != HttpStatus.SC_NO_CONTENT) {
                        //System.out.println(delResp.getHttpBody());
                    }
                }
            }
            catch(Exception e) {
                throw new RuntimeException(e);
            }
        }
    }

    protected Email searchMails(User user, String mailType) {
        List<Email> emails = user.getEmails();
        for(Email m : emails) {
            if(m.getType().equalsIgnoreCase(mailType)) {
                return m;
            }
        }
        
        return null;
    }
    
    protected User buildUser() {
        String username = randomAlphabetic(5);
        User user = new User();
        user.setUserName(username);

        Name name = new Name();
        name.setFamilyName(username);
        name.setGivenName(username);
        name.setHonorificPrefix("Mr.");
        name.setFormatted(name.getHonorificPrefix() + " " + name.getGivenName() + " " + name.getFamilyName());
        user.setName(name);
        
        List<Email> emails = new ArrayList<Email>();
        
        Email homeMail = new Email();
        homeMail.setDisplay("Home Email");
        homeMail.setType("home");
        String s = randomAlphabetic(5);
        homeMail.setValue(s + "@home.com" );
        emails.add(homeMail);
        
        Email workMail = new Email();
        workMail.setDisplay("Work Email");
        workMail.setType("work");
        s = randomAlphabetic(5);
        workMail.setValue(s + "@work.com" );
        emails.add(workMail);
        
        user.setEmails(emails);
        user.setActive(true);
        
        user.setPassword(randomAlphabetic(11));
        
        return user;
    }
}
