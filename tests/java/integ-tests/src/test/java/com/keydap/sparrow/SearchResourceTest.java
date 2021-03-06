/*
 * Copyright (c) 2016 Keydap Software.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * See LICENSE file for details.
 */
package com.keydap.sparrow;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.fail;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import org.apache.http.HttpStatus;
import org.junit.BeforeClass;
import org.junit.Test;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.keydap.sparrow.scim.Device;
import com.keydap.sparrow.scim.Group;
import com.keydap.sparrow.scim.User;
import com.keydap.sparrow.scim.User.Address;
import com.keydap.sparrow.scim.User.Email;
import com.keydap.sparrow.scim.User.EnterpriseUser;
import com.keydap.sparrow.scim.User.Name;

/**
 *
 * @author Kiran Ayyagari (kayyagari@keydap.com)
 */
public class SearchResourceTest extends TestBase {
    
    private static User admin;
    private static User snowden;
    private static User bhagat;
    private static User assange;
    private static User stallman;
    
    private static Device thermostat;
    private static Device watch;
    private static Device mobile;
    
    private static String tDate = "1947-08-14T18:30:00Z";
    private static String wDate = "2016-05-04T14:19:14Z";
    private static String mDate = "2015-09-19T01:30:00Z";
    
    @BeforeClass
    public static void cleanAndInject() throws Exception {
        deleteAll(User.class);
        deleteAll(Device.class);
        deleteAll(Group.class);
        deleteAll(RegisteredApp.class);

        // 'admin' user is the existing user, hence NOT inserted 
        // from here. It is present only to help some tests pass
        admin = new User();
        admin.setUserName("admin");
        admin.setActive(true);
        admin.setDisplayName("Administrator");
        
        Email adEmail = new Email();
        adEmail.setValue("admin@example.com");
        adEmail.setType("work");
        adEmail.setPrimary(true);
        admin.setEmails(Collections.singletonList(adEmail));

        snowden = new User();
        snowden.setUserName("snowden");
        snowden.setActive(true);
        snowden.setDisplayName("Edward J. Snowden");
        Name sName = new Name();
        sName.setFamilyName("Edward");
        sName.setFormatted(snowden.getDisplayName());
        sName.setGivenName("Snowden");
        sName.setHonorificPrefix("Mr.");
        snowden.setName(sName);
        
        Email sEmail1 = new Email();
        sEmail1.setValue("sn@eff.org");
        sEmail1.setType("work");
        sEmail1.setPrimary(true);
        
        Email sEmail2 = new Email();
        sEmail2.setValue("sn@snowden.com");
        sEmail2.setType("home");
        sEmail2.setPrimary(false);

        List<Email> sEmails = new ArrayList<Email>();
        sEmails.add(sEmail1);
        sEmails.add(sEmail2);
        snowden.setEmails(sEmails);
        
        Address sAddress = new Address();
        sAddress.setCountry("RU");
        sAddress.setLocality("St. Petersburg");
        sAddress.setStreetAddress("1st Avenue");
        snowden.setAddresses(Collections.singletonList(sAddress));
        
        snowden.setPassword("Secret001");
        snowden.setNickName("hero");

        assange = new User();
        assange.setUserName("assange");
        assange.setActive(true);
        assange.setDisplayName("Julien Assange");
        Name aName = new Name();
        aName.setFamilyName("Julien");
        aName.setFormatted(assange.getDisplayName());
        aName.setGivenName("Assange");
        aName.setHonorificPrefix("Mr.");
        assange.setName(aName);
        
        Email aEmail = new Email();
        aEmail.setDisplay("");
        aEmail.setValue("assange@wikileaks.org");
        aEmail.setType("home");
        aEmail.setPrimary(true);
        assange.setEmails(Collections.singletonList(aEmail));
        
        Address aAddress = new Address();
        aAddress.setCountry("UK");
        aAddress.setLocality("Ecaudor Embassy");
        aAddress.setStreetAddress("1st Avenue");
        assange.setAddresses(Collections.singletonList(aAddress));
        
        assange.setPassword("Secret002");
        assange.setNickName("pioneer");
        
        bhagat = new User();
        bhagat.setUserName("bhagat");
        bhagat.setActive(true);
        bhagat.setDisplayName("Bhagat Singh");
        Name bName = new Name();
        bName.setFamilyName("Singh");
        bName.setFormatted(bhagat.getDisplayName());
        bName.setGivenName("Bhagat");
        bName.setHonorificPrefix("Mr.");
        bhagat.setName(bName);
        
        Email bEmail = new Email();
        bEmail.setDisplay("");
        bEmail.setValue("bhagat@hra.org");
        bEmail.setType("home");
        bEmail.setPrimary(true);
        bhagat.setEmails(Collections.singletonList(bEmail));
        
        Address bAddress = new Address();
        bAddress.setCountry("IN");
        bAddress.setLocality("Punjab");
        bAddress.setStreetAddress("2nd Avenue");
        bhagat.setAddresses(Collections.singletonList(bAddress));
        
        bhagat.setPassword("Secret003");
        bhagat.setNickName("martyr");
        
        stallman = new User();
        stallman.setUserName("stallman");
        stallman.setActive(true);
        stallman.setDisplayName("Richard M. Stallman");
        Name rName = new Name();
        rName.setFamilyName("Richard");
        rName.setFormatted(stallman.getDisplayName());
        rName.setGivenName("Stallman");
        rName.setHonorificPrefix("Mr.");
        stallman.setName(rName);
        
        Email rEmail = new Email();
        rEmail.setDisplay("");
        rEmail.setValue("rms@fsf.org");
        rEmail.setType("home");
        rEmail.setPrimary(true);
        stallman.setEmails(Collections.singletonList(rEmail));
        
        Address rAddress = new Address();
        rAddress.setCountry("US");
        rAddress.setLocality("New York City");
        rAddress.setStreetAddress("2nd Avenue");
        stallman.setAddresses(Collections.singletonList(rAddress));
        
        stallman.setPassword("Secret004");
        stallman.setNickName("pioneer");

        EnterpriseUser stallmanEu = new EnterpriseUser();
        stallmanEu.setCostCenter("GCC");
        stallmanEu.setDivision("GNU");
        stallmanEu.setEmployeeNumber("1");
        stallmanEu.setOrganization("FSF");
        stallman.setEnterpriseUser(stallmanEu);

        client.addResource(snowden);
        client.addResource(assange);
        client.addResource(bhagat);
        client.addResource(stallman);
        
        thermostat = new Device();
        thermostat.setManufacturer("Samsung");
        thermostat.setPrice(900.07);
        thermostat.setRating(2);
        thermostat.setSerialNumber("011");
        thermostat.setInstalledDate(utcDf.parse(tDate));
        
        watch = new Device();
        watch.setManufacturer("Fossil");
        watch.setPrice(8000);
        watch.setRating(9);
        watch.setSerialNumber("002");
        watch.setInstalledDate(utcDf.parse(wDate));
        
        mobile = new  Device();
        mobile.setManufacturer("Apple");
        mobile.setPrice(53500);
        mobile.setRating(10);
        mobile.setSerialNumber("007");
        mobile.setInstalledDate(utcDf.parse(mDate));
        
        client.addResource(thermostat);
        client.addResource(watch);
        client.addResource(mobile);
    }
    
    @Test
    public void testStringComparison() {
        SearchResponse<User> resp = client.searchResource("username eq \"snowden\"", User.class);
        checkResults(resp, snowden);
        User found = resp.getResources().get(0);
        assertNull(found.getPassword()); // password should never be returned
        
        resp = client.searchResource("emails.type eq \"work\"", User.class);
        checkResults(resp, snowden, admin);
        
        resp = client.searchResource("username ne \"snowden\"", User.class);
        checkResults(resp, assange, bhagat, stallman, admin);

        resp = client.searchResource("name.formatted co \"l\"", User.class);
        checkResults(resp, stallman, assange);

        resp = client.searchResource("name.formatted co \"L\"", User.class);
        checkResults(resp, stallman, assange);

        resp = client.searchResource("name.familyName sw \"J\"", User.class);
        checkResults(resp, assange);

        resp = client.searchResource("name.familyName sw \"j\"", User.class);
        checkResults(resp, assange);

        resp = client.searchResource("emails.value ew \".org\"", User.class);
        checkResults(resp, stallman, assange, snowden, bhagat);

        resp = client.searchResource("emails.value ew \".com\"", User.class);
        checkResults(resp, snowden, admin);
        
        resp = client.searchResource("emails.value ew \".COM\"", User.class);
        checkResults(resp, snowden, admin);
        
        resp = client.searchResource("costCenter pr", User.class);
        checkResults(resp, stallman);
    }

    @Test
    public void testArithmeticOperators() {
        SearchResponse<User> uresp = client.searchResource("username gt \"snowden\"", User.class);
        checkResults(uresp, stallman);

        uresp = client.searchResource("not username gt \"snowden\"", User.class);
        checkResults(uresp, assange, bhagat, snowden, admin);
        
        uresp = client.searchResource("username ge \"snowden\"", User.class);
        checkResults(uresp, stallman, snowden);
        
        uresp = client.searchResource("not username ge \"snowden\"", User.class);
        checkResults(uresp, assange, bhagat, admin);
        
        uresp = client.searchResource("username lt \"snowden\"", User.class);
        checkResults(uresp, assange, bhagat, admin);  

        uresp = client.searchResource("username le \"snowden\"", User.class);
        checkResults(uresp, assange, snowden, bhagat, admin);        

        SearchResponse<Device> dresp = client.searchResource("rating gt  ", Device.class);
        assertEquals(HttpStatus.SC_BAD_REQUEST, dresp.getHttpCode());
        
        dresp = client.searchResource("rating eq 9", Device.class);
        checkResults(dresp, watch);

        dresp = client.searchResource("rating gt  9", Device.class);
        checkResults(dresp, mobile);

        dresp = client.searchResource("rating ge 9", Device.class);
        checkResults(dresp, watch, mobile);

        dresp = client.searchResource("price lt 900.10", Device.class);
        checkResults(dresp, thermostat);
        
        dresp = client.searchResource("price le 8000.10", Device.class);
        checkResults(dresp, thermostat, watch);

        dresp = client.searchResource("installedDate le \"2016-05-04T14:19:14Z\"", Device.class);
        checkResults(dresp, thermostat, watch, mobile);
    }
        
    @Test
    public void testLogicalOperators() {
        SearchResponse<User> resp = client.searchResource("emails[type eq \"work\" or value co \"org\"]", User.class);
        checkResults(resp, snowden, assange, bhagat, stallman, admin);

        resp = client.searchResource("emails[type eq \"work\" and value co \"org\"]", User.class);
        checkResults(resp, snowden);

        resp = client.searchResource("emails.type eq \"work\" and emails.value co \"com\"", User.class);
        checkResults(resp, snowden, admin);
        
        resp = client.searchResource("unknownAttribute eq \"work\" and emails.value co \"com\"", User.class);
        checkResults(resp);
        
        resp = client.searchResource("not costCenter pr", User.class);
        checkResults(resp, snowden, assange, bhagat, admin);
    }

    @Test
    public void testFilterWithNestedOperators() {
        SearchResponse<User> resp = client.searchResource("not emails[type eq \"work\" or value co \"org\"]", User.class);
        checkResults(resp);

        resp = client.searchResource("not emails[type eq \"work\"] and username co \"ss\"", User.class);
        checkResults(resp, assange, snowden, bhagat, stallman, admin);

        resp = client.searchResource("(not emails[type eq \"work\"]) and username co \"ss\"", User.class);
        checkResults(resp, assange);
        
        resp = client.searchResource("schemas eq \"" + User.SCHEMA + "\"", User.class);
        checkResults(resp, assange, snowden, bhagat, stallman, admin);

        // emails.value is indexed
        resp = client.searchResource("emails.value co \"org\" and (username co \"ss\" and displayname sw \"j\")", User.class);
        checkResults(resp, assange);

        // now try to fetch the same result but using a non indexed emails sub-attribute
        resp = client.searchResource("(emails.type co \"home\" and (((username co \"ss\" ))))and displayname sw \"j\"", User.class);
        checkResults(resp, assange);
    }

    @Test
    public void testWithSearchRequest() {
        SearchRequest req = new SearchRequest();
        req.setFilter("emails.type eq \"work\"");
        req.setAttributes("username");
        
        SearchResponse<User> resp = client.searchResource(req, User.class);
        checkResults(resp, snowden, admin);
        // only the ID, schemas and username fields should be present
        User fetched = resp.getResources().get(0);
        assertNull(fetched.getEmails());
        assertNull(fetched.getDisplayName());
        
        req = new SearchRequest();
        req.setFilter("emails.type eq \"work\"");
        req.setExcludedAttributes("emails");
        
        resp = client.searchResource(req, User.class);
        checkResults(resp, snowden, admin);
        // only the emails should NOT present
        fetched = resp.getResources().get(0);
        assertNull(fetched.getEmails());
        assertNotNull(fetched.getDisplayName());
        
        req = new SearchRequest();
        // use the pr operator on an indexed field
        req.setFilter("username pr");
        resp = client.searchResource(req, User.class);
        checkResults(resp, snowden, stallman, bhagat, assange, admin);

        req = new SearchRequest();
        // use the pr operator on an indexed field
        req.setFilter("emails[type eq \"home\" and name[formatted pr]");
        resp = client.searchResource(req, User.class);
        assertEquals(HttpStatus.SC_BAD_REQUEST, resp.getHttpCode());
        assertNotNull(resp.getError());

        req = new SearchRequest();
        // use the pr operator on an indexed field
        req.setFilter("name[formatted pr]");
        resp = client.searchResource(req, User.class);
        checkResults(resp, snowden, stallman, bhagat, assange);
    }
    
    @Test
    public void testAtRoot() {
        SearchRequest req = new SearchRequest();
        req.setFilter("id pr");

        SearchResponse<Object> resp = client.searchAll(req);
        assertEquals(HttpStatus.SC_OK, resp.getHttpCode());
        List<Object> received = resp.getResources();
        assertEquals(10, received.size());
    }

    @Test
    public void testFilterWithPrefixedAts() {
        SearchResponse<User> uresp = client.searchResource(EnterpriseUser.SCHEMA + ":employeeNumber pr", User.class);
        checkResults(uresp, stallman);
        JsonObject json = (JsonObject) new JsonParser().parse(uresp.getHttpBody());
        JsonObject rs = (JsonObject) json.get("Resources").getAsJsonArray().get(0);
        assertNotNull(rs.get(EnterpriseUser.SCHEMA));

        uresp = client.searchResource(EnterpriseUser.SCHEMA + ":employeeNumber pr", User.class, "employeeNumber", "username");
        checkResults(uresp, stallman);
        json = (JsonObject) new JsonParser().parse(uresp.getHttpBody());
        rs = (JsonObject) json.get("Resources").getAsJsonArray().get(0);
        assertNotNull(rs.get(EnterpriseUser.SCHEMA));

        uresp = client.searchResource(User.SCHEMA.toLowerCase() + ":emails pr", User.class, "employeeNumber", "username");
        checkResults(uresp, stallman, assange, snowden, bhagat, admin);
    }
    
    private void checkResults(SearchResponse<User> resp, User... ids) {
        assertEquals(HttpStatus.SC_OK, resp.getHttpCode());
        List<User> received = resp.getResources();
        
        if(ids != null && ids.length > 0) {
            int expectedCount = ids.length;
            assertEquals(expectedCount, received.size());
            for(User r : received) {
                for(User i : ids) {
                    if (i.getUserName().equals(r.getUserName())) {
                        expectedCount--;
                    }
                }
            }
            
            if(expectedCount != 0) {
                fail("All the expected users are not present in the Response");
            }
        } else {
            assertNull(received);
        }
    }

    private void checkResults(SearchResponse<Device> resp, Device... ids) {
        assertEquals(HttpStatus.SC_OK, resp.getHttpCode());
        List<Device> received = resp.getResources();
        
        if(ids != null && ids.length > 0) {
            int expectedCount = ids.length;
            assertEquals(expectedCount, received.size());
            for(Device r : received) {
                for(Device i : ids) {
                    if (i.getSerialNumber().equals(r.getSerialNumber())) {
                        expectedCount--;
                    }
                }
            }
            
            if(expectedCount != 0) {
                fail("All the expected devices are not present in the Response");
            }
        } else {
            assertNull(received);
        }
    }
}
