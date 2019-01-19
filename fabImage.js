/*
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
*/

'use strict';
const shim = require('fabric-shim');
const util = require('util');

let Chaincode = class {

  // The Init method is called when the Smart Contract 'fabImage' is instantiated by the blockchain network
  // Best practice is to have any Ledger initialization in separate function -- see initLedger()
  async Init(stub) {
    console.info('=========== Instantiated fabImage chaincode ===========');
    return shim.success();
  }

  // The Invoke method is called as a result of an application request to run the Smart Contract
  // 'fabcar'. The calling application program has also specified the particular smart contract
  // function to be called, with arguments
  async Invoke(stub) {
    let ret = stub.getFunctionAndParameters();
    console.info(ret);

    let method = this[ret.fcn];
    if (!method) {
      console.error('no function of name:' + ret.fcn + ' found');
      throw new Error('Received unknown function ' + ret.fcn + ' invocation');
    }
    try {
      let payload = await method(stub, ret.params);
      return shim.success(payload);
    } catch (err) {
      console.log(err);
      return shim.error(err);
    }
  }

  async queryImage(stub, args) {
    if (args.length != 1) {
      throw new Error('Incorrect number of arguments. Expecting ImageId ex: IMG01');
    }
    let ImageId = args[0];

    let imageAsBytes = await stub.getState(ImageId); //get the car from chaincode state
    if (!imageAsBytes || imageAsBytes.toString().length <= 0) {
      throw new Error(ImageId + ' does not exist: ');
    }
    console.log(imageAsBytes.toString());
    return imageAsBytes;
  }

  async initLedger(stub, args) {
    console.info('============= START : Initialize Ledger ===========');
    let imgs = [];
    imgs.push({
      imageName: 'Typhoid',
      imageSize: ' 1 MB',
      Owner: 'Tomoko'
    });
    imgs.push({
      imageName: 'Pnemonia',
      imageSize: ' 1 MB',
      Owner: 'Jin'
    });
    imgs.push({
      imageName: 'Anemia',
      imageSize: ' 1 MB',
      Owner: 'Max'
    });
   

    for (let i = 0; i < imgs.length; i++) {
      imgs[i].docType = 'img';
      await stub.putState('IMG' + i, Buffer.from(JSON.stringify(imgs[i])));
      console.info('Added <--> ', imgs[i]);
    }
    console.info('============= END : Initialize Ledger ===========');
  }

  async UploadImage(stub, args) {
    console.info('============= START : Upload Image ===========');
    if (args.length != 3) {
      throw new Error('Incorrect number of arguments. Expecting 3');
    }

    var img = {
      docType: 'img',
      imageName: args[1],
      imageSize: args[2],
      Owner: args[3]
    };

    await stub.putState(args[0], Buffer.from(JSON.stringify(img)));
    console.info('============= END : Upload Image ===========');
  }

  async queryAllImgs(stub, args) {

    let startKey = 'IMG0';
    let endKey = 'IMG999';

    let iterator = await stub.getStateByRange(startKey, endKey);

    let allResults = [];
    while (true) {
      let res = await iterator.next();

      if (res.value && res.value.value.toString()) {
        let jsonRes = {};
        console.log(res.value.value.toString('utf8'));

        jsonRes.Key = res.value.key;
        try {
          jsonRes.Record = JSON.parse(res.value.value.toString('utf8'));
        } catch (err) {
          console.log(err);
          jsonRes.Record = res.value.value.toString('utf8');
        }
        allResults.push(jsonRes);
      }
      if (res.done) {
        console.log('end of data');
        await iterator.close();
        console.info(allResults);
        return Buffer.from(JSON.stringify(allResults));
      }
    }
  }

  async transferImage(stub, args) {
    console.info('============= START : transferImage ===========');
    if (args.length != 3) {
      throw new Error('Incorrect number of arguments. Expecting 3');
    }

    let imageAsBytes = await stub.getState(args[0]);
    let img = JSON.parse(imageAsBytes);
    img.Owner = args[1];
    img.imageName = args[2];

    await stub.putState(args[0], Buffer.from(JSON.stringify(img)));
    console.info('============= END : transferImage ===========');
  }
};

shim.start(new Chaincode());
