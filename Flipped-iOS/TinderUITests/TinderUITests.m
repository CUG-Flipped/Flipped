//
//  TinderUITests.m
//  TinderUITests
//
//  Created by Layer on 2020/4/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <XCTest/XCTest.h>

@interface TinderUITests : XCTestCase

@end

@implementation TinderUITests

- (void)setUp {
    // Put setup code here. This method is called before the invocation of each test method in the class.

    // In UI tests it is usually best to stop immediately when a failure occurs.
    self.continueAfterFailure = NO;

    // In UI tests it’s important to set the initial state - such as interface orientation - required for your tests before they run. The setUp method is a good place to do this.
}

- (void)tearDown {
    // Put teardown code here. This method is called after the invocation of each test method in the class.
}

#pragma mark -- 注册
- (void)testReginstration {
    XCUIApplication *app = [[XCUIApplication alloc] init];
    [app launch];
    
    XCUIElement *enterFullNameTextField = app.textFields[@"Enter full name"];
    [enterFullNameTextField tap];
    
    XCUIElement *tKey = app/*@START_MENU_TOKEN@*/.keys[@"T"]/*[[".keyboards.keys[@\"T\"]",".keys[@\"T\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [tKey tap];
    
    XCUIElement *eKey = app/*@START_MENU_TOKEN@*/.keys[@"e"]/*[[".keyboards.keys[@\"e\"]",".keys[@\"e\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [eKey tap];
    
    XCUIElement *sKey = app/*@START_MENU_TOKEN@*/.keys[@"s"]/*[[".keyboards.keys[@\"s\"]",".keys[@\"s\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [sKey tap];

    XCUIElement *tKey2 = app/*@START_MENU_TOKEN@*/.keys[@"t"]/*[[".keyboards.keys[@\"t\"]",".keys[@\"t\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [tKey2 tap];
    
    XCUIElement *element = [[[[[[[[[app childrenMatchingType:XCUIElementTypeWindow] elementBoundByIndex:0] childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element;
    [element tap];
    [app.textFields[@"Enter email"] tap];
    [tKey tap];
    [eKey tap];
    [sKey tap];
    [tKey2 tap];
    
    XCUIElement *moreKey = app/*@START_MENU_TOKEN@*/.keys[@"more"]/*[[".keyboards",".keys[@\"numbers\"]",".keys[@\"more\"]"],[[[-1,2],[-1,1],[-1,0,1]],[[-1,2],[-1,1]]],[0]]@END_MENU_TOKEN@*/;
    [moreKey tap];
    
    XCUIElement *key = app/*@START_MENU_TOKEN@*/.keys[@"@"]/*[[".keyboards.keys[@\"@\"]",".keys[@\"@\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key tap];
    
    XCUIElement *key2 = app/*@START_MENU_TOKEN@*/.keys[@"1"]/*[[".keyboards.keys[@\"1\"]",".keys[@\"1\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key2 tap];
    
    XCUIElement *key3 = app/*@START_MENU_TOKEN@*/.keys[@"6"]/*[[".keyboards.keys[@\"6\"]",".keys[@\"6\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key3 tap];
    
    XCUIElement *key4 = app/*@START_MENU_TOKEN@*/.keys[@"3"]/*[[".keyboards.keys[@\"3\"]",".keys[@\"3\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key4 tap];
    
    XCUIElement *key5 = app/*@START_MENU_TOKEN@*/.keys[@"."]/*[[".keyboards.keys[@\".\"]",".keys[@\".\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key5 tap];
    
    XCUIElement *moreKey2 = app/*@START_MENU_TOKEN@*/.keys[@"more"]/*[[".keyboards",".keys[@\"letters\"]",".keys[@\"more\"]"],[[[-1,2],[-1,1],[-1,0,1]],[[-1,2],[-1,1]]],[0]]@END_MENU_TOKEN@*/;
    [moreKey2 tap];
    
    XCUIElement *cKey = app/*@START_MENU_TOKEN@*/.keys[@"c"]/*[[".keyboards.keys[@\"c\"]",".keys[@\"c\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [cKey tap];
    
    XCUIElement *oKey = app/*@START_MENU_TOKEN@*/.keys[@"o"]/*[[".keyboards.keys[@\"o\"]",".keys[@\"o\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [oKey tap];
    
    XCUIElement *mKey = app/*@START_MENU_TOKEN@*/.keys[@"m"]/*[[".keyboards.keys[@\"m\"]",".keys[@\"m\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [mKey tap];
    [element tap];
    [app.secureTextFields[@"Enter password"] tap];
    
    XCUIElement *qKey = app/*@START_MENU_TOKEN@*/.keys[@"q"]/*[[".keyboards.keys[@\"q\"]",".keys[@\"q\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [qKey tap];
    
    XCUIElement *wKey = app/*@START_MENU_TOKEN@*/.keys[@"w"]/*[[".keyboards.keys[@\"w\"]",".keys[@\"w\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [wKey tap];
    [eKey tap];
    
    XCUIElement *rKey = app/*@START_MENU_TOKEN@*/.keys[@"r"]/*[[".keyboards.keys[@\"r\"]",".keys[@\"r\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [rKey tap];
    [tKey2 tap];
    
    XCUIElement *yKey = app/*@START_MENU_TOKEN@*/.keys[@"y"]/*[[".keyboards.keys[@\"y\"]",".keys[@\"y\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [yKey tap];
    [element tap];
    
    XCUIElement *registerButton = app.buttons[@"Register"];
    [registerButton tap];
    [enterFullNameTextField tap];
    [wKey tap];
    [wKey tap];
    [wKey tap];
    [wKey tap];
    [wKey tap];
    [element tap];
    [registerButton tap];
}

#pragma mark -- 登录
- (void)testLogin {
    XCUIApplication *app = [[XCUIApplication alloc] init];
    [app launch];
    
    [app.textFields[@"Enter email"] tap];
    
    XCUIElement *tKey = app/*@START_MENU_TOKEN@*/.keys[@"T"]/*[[".keyboards.keys[@\"T\"]",".keys[@\"T\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [tKey tap];
    
    XCUIElement *eKey = app/*@START_MENU_TOKEN@*/.keys[@"e"]/*[[".keyboards.keys[@\"e\"]",".keys[@\"e\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [eKey tap];
    
    XCUIElement *sKey = app/*@START_MENU_TOKEN@*/.keys[@"s"]/*[[".keyboards.keys[@\"s\"]",".keys[@\"s\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [sKey tap];
    
    XCUIElement *tKey2 = app/*@START_MENU_TOKEN@*/.keys[@"t"]/*[[".keyboards.keys[@\"t\"]",".keys[@\"t\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [tKey2 tap];
    
    XCUIElement *moreKey = app/*@START_MENU_TOKEN@*/.keys[@"more"]/*[[".keyboards",".keys[@\"numbers\"]",".keys[@\"more\"]"],[[[-1,2],[-1,1],[-1,0,1]],[[-1,2],[-1,1]]],[0]]@END_MENU_TOKEN@*/;
    [moreKey tap];
    
    XCUIElement *key = app/*@START_MENU_TOKEN@*/.keys[@"@"]/*[[".keyboards.keys[@\"@\"]",".keys[@\"@\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key tap];
    
    XCUIElement *key2 = app/*@START_MENU_TOKEN@*/.keys[@"1"]/*[[".keyboards.keys[@\"1\"]",".keys[@\"1\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key2 tap];
    
    XCUIElement *key3 = app/*@START_MENU_TOKEN@*/.keys[@"6"]/*[[".keyboards.keys[@\"6\"]",".keys[@\"6\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key3 tap];
    
    XCUIElement *key4 = app/*@START_MENU_TOKEN@*/.keys[@"3"]/*[[".keyboards.keys[@\"3\"]",".keys[@\"3\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key4 tap];
    
    XCUIElement *key5 = app/*@START_MENU_TOKEN@*/.keys[@"."]/*[[".keyboards.keys[@\".\"]",".keys[@\".\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key5 tap];
    
    XCUIElement *moreKey2 = app/*@START_MENU_TOKEN@*/.keys[@"more"]/*[[".keyboards",".keys[@\"letters\"]",".keys[@\"more\"]"],[[[-1,2],[-1,1],[-1,0,1]],[[-1,2],[-1,1]]],[0]]@END_MENU_TOKEN@*/;
    [moreKey2 tap];
    
    XCUIElement *cKey = app/*@START_MENU_TOKEN@*/.keys[@"c"]/*[[".keyboards.keys[@\"c\"]",".keys[@\"c\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [cKey tap];
    
    XCUIElement *oKey = app/*@START_MENU_TOKEN@*/.keys[@"o"]/*[[".keyboards.keys[@\"o\"]",".keys[@\"o\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [oKey tap];
    
    XCUIElement *mKey = app/*@START_MENU_TOKEN@*/.keys[@"m"]/*[[".keyboards.keys[@\"m\"]",".keys[@\"m\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [mKey tap];
    
    XCUIElement *element = [[[[[[[[[app childrenMatchingType:XCUIElementTypeWindow] elementBoundByIndex:0] childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element;
    [element tap];
    [app.secureTextFields[@"Enter password"] tap];
    
    XCUIElement *qKey = app/*@START_MENU_TOKEN@*/.keys[@"q"]/*[[".keyboards.keys[@\"q\"]",".keys[@\"q\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [qKey tap];
    [app/*@START_MENU_TOKEN@*/.keys[@"w"]/*[[".keyboards.keys[@\"w\"]",".keys[@\"w\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    [eKey tap];
    
    XCUIElement *rKey = app/*@START_MENU_TOKEN@*/.keys[@"r"]/*[[".keyboards.keys[@\"r\"]",".keys[@\"r\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [rKey tap];
    [tKey2 tap];
    
    XCUIElement *yKey = app/*@START_MENU_TOKEN@*/.keys[@"y"]/*[[".keyboards.keys[@\"y\"]",".keys[@\"y\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [yKey tap];
    [element tap];
    [app.buttons[@"Log In"] tap];
}

#pragma mark -- 首页
- (void)testHome {
    
    XCUIApplication *app = [[XCUIApplication alloc] init];
    [app launch];
    
    XCUIElement *element = [[[[[[[[[[app childrenMatchingType:XCUIElementTypeWindow] elementBoundByIndex:0] childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther] elementBoundByIndex:1];
    [element swipeLeft];
    
    [element swipeRight];
    
    [element swipeUp];
    
    [element swipeDown];
    
    [app.buttons[@"twoImg"] tap];
    
    [app.buttons[@"fourImg"] tap];
    
    [app.buttons[@"threeImg"] tap];
    
    [app.buttons[@"firstImg"] tap];
    
    [app.buttons[@"fiveImg"] tap];
    
}

#pragma mark -- 个人信息
- (void)testInfo1 {
    
    XCUIApplication *app = [[XCUIApplication alloc] init];
    [app launch];
    
    [app.buttons[@"navLeftImg"] tap];
    
    XCUIElementQuery *tablesQuery2 = app.tables;
    [[[[tablesQuery2 childrenMatchingType:XCUIElementTypeCell] elementBoundByIndex:0] childrenMatchingType:XCUIElementTypeTextField].element tap];
    
    XCUIElementQuery *tablesQuery = tablesQuery2;
    [tablesQuery/*@START_MENU_TOKEN@*/.buttons[@"Clear text"]/*[[".cells",".textFields.buttons[@\"Clear text\"]",".buttons[@\"Clear text\"]"],[[[-1,2],[-1,1],[-1,0,1]],[[-1,2],[-1,1]]],[0]]@END_MENU_TOKEN@*/ tap];
    
    XCUIElement *lKey = app/*@START_MENU_TOKEN@*/.keys[@"L"]/*[[".keyboards.keys[@\"L\"]",".keys[@\"L\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [lKey tap];
    
    XCUIElement *aKey = app/*@START_MENU_TOKEN@*/.keys[@"a"]/*[[".keyboards.keys[@\"a\"]",".keys[@\"a\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [aKey tap];
    
    XCUIElement *yKey = app/*@START_MENU_TOKEN@*/.keys[@"y"]/*[[".keyboards.keys[@\"y\"]",".keys[@\"y\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [yKey tap];
    
    [app/*@START_MENU_TOKEN@*/.keys[@"e"]/*[[".keyboards.keys[@\"e\"]",".keys[@\"e\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    
    XCUIElement *rKey = app/*@START_MENU_TOKEN@*/.keys[@"r"]/*[[".keyboards.keys[@\"r\"]",".keys[@\"r\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [rKey tap];
    
    [tablesQuery/*@START_MENU_TOKEN@*/.staticTexts[@"Profession"]/*[[".otherElements[@\"Profession\"].staticTexts[@\"Profession\"]",".staticTexts[@\"Profession\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ swipeUp];
    [app/*@START_MENU_TOKEN@*/.staticTexts[@"Save"]/*[[".buttons[@\"Save\"].staticTexts[@\"Save\"]",".staticTexts[@\"Save\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    [app.buttons[@"Cancel"] tap];
    
    
}

#pragma mark -- 详细信息
- (void)testInfo2 {
    
    XCUIApplication *app = [[XCUIApplication alloc] init];
    
    [app launch];
    
    [[[[[[[[[[[app childrenMatchingType:XCUIElementTypeWindow] elementBoundByIndex:0] childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther] elementBoundByIndex:1].buttons[@"i"] tap];
    
    [app.buttons[@"inforNoLike"] tap];
    
    [app.buttons[@"inforLike"] tap];
    
    [app.buttons[@"inforSC"] tap];
    
    [app.buttons[@"infoBackDown"] tap];
    
}

#pragma mark -- 通讯录及聊天
- (void)testChat {
    
    XCUIApplication *app = [[XCUIApplication alloc] init];
    [app.buttons[@"navRightImg"] tap];
    
    [app.tables/*@START_MENU_TOKEN@*/.staticTexts[@"Hermosa"]/*[[".cells.staticTexts[@\"Hermosa\"]",".staticTexts[@\"Hermosa\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    [app.staticTexts[@"Please enter here..."] tap];
    
    XCUIElement *hKey = app/*@START_MENU_TOKEN@*/.keys[@"H"]/*[[".keyboards.keys[@\"H\"]",".keys[@\"H\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [hKey tap];
    
    XCUIElement *eKey = app/*@START_MENU_TOKEN@*/.keys[@"e"]/*[[".keyboards.keys[@\"e\"]",".keys[@\"e\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [eKey tap];
    
    XCUIElement *lKey = app/*@START_MENU_TOKEN@*/.keys[@"l"]/*[[".keyboards.keys[@\"l\"]",".keys[@\"l\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [lKey tap];
    
    XCUIElement *oKey = app/*@START_MENU_TOKEN@*/.keys[@"o"]/*[[".keyboards.keys[@\"o\"]",".keys[@\"o\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [oKey tap];
    [oKey tap];
    [app.buttons[@"btn expression"] tap];
//    [[app.scrollViews containingType:XCUIElementTypeStaticText identifier:@"\U4eba\U7269"].element tap];
    [[[[[[[[[[app childrenMatchingType:XCUIElementTypeWindow] elementBoundByIndex:0] childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeTable].element tap];
    
    
}


- (void)testLaunchPerformance {
    if (@available(macOS 10.15, iOS 13.0, tvOS 13.0, *)) {
        // This measures how long it takes to launch your application.
        [self measureWithMetrics:@[XCTOSSignpostMetric.applicationLaunchMetric] block:^{
            [[[XCUIApplication alloc] init] launch];
        }];
    }
}
@end
