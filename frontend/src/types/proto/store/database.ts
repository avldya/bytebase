/* eslint-disable */
import * as Long from "long";
import * as _m0 from "protobufjs/minimal";
import { StringValue } from "../google/protobuf/wrappers";

export const protobufPackage = "bytebase.store";

/** DatabaseMetadata is the metadata for databases. */
export interface DatabaseMetadata {
  name: string;
  /** The schemas is the list of schemas in a database. */
  schemas: SchemaMetadata[];
  /** The character_set is the character set of a database. */
  characterSet: string;
  /** The collation is the collation of a database. */
  collation: string;
  /** The extensions is the list of extensions in a database. */
  extensions: ExtensionMetadata[];
}

/**
 * SchemaMetadata is the metadata for schemas.
 * This is the concept of schema in Postgres, but it's a no-op for MySQL.
 */
export interface SchemaMetadata {
  /**
   * The name is the schema name.
   * It is an empty string for databases without such concept such as MySQL.
   */
  name: string;
  /** The tables is the list of tables in a schema. */
  tables: TableMetadata[];
  /** The views is the list of views in a schema. */
  views: ViewMetadata[];
  /** The functions is the list of functions in a schema. */
  functions: FunctionMetadata[];
}

/** TableMetadata is the metadata for tables. */
export interface TableMetadata {
  /** The name is the name of a table. */
  name: string;
  /** The columns is the ordered list of columns in a table. */
  columns: ColumnMetadata[];
  /** The indexes is the list of indexes in a table. */
  indexes: IndexMetadata[];
  /** The engine is the engine of a table. */
  engine: string;
  /** The collation is the collation of a table. */
  collation: string;
  /** The row_count is the estimated number of rows of a table. */
  rowCount: number;
  /** The data_size is the estimated data size of a table. */
  dataSize: number;
  /** The index_size is the estimated index size of a table. */
  indexSize: number;
  /** The data_free is the estimated free data size of a table. */
  dataFree: number;
  /** The create_options is the create option of a table. */
  createOptions: string;
  /** The comment is the comment of a table. */
  comment: string;
  /** The foreign_keys is the list of foreign keys in a table. */
  foreignKeys: ForeignKeyMetadata[];
}

/** ColumnMetadata is the metadata for columns. */
export interface ColumnMetadata {
  /** The name is the name of a column. */
  name: string;
  /** The position is the position in columns. */
  position: number;
  /** The default is the default of a column. Use google.protobuf.StringValue to distinguish between an empty string default value or no default. */
  default?: string;
  /** The nullable is the nullable of a column. */
  nullable: boolean;
  /** The type is the type of a column. */
  type: string;
  /** The character_set is the character_set of a column. */
  characterSet: string;
  /** The collation is the collation of a column. */
  collation: string;
  /** The comment is the comment of a column. */
  comment: string;
}

/** ViewMetadata is the metadata for views. */
export interface ViewMetadata {
  /** The name is the name of a view. */
  name: string;
  /** The definition is the definition of a view. */
  definition: string;
  /** The comment is the comment of a view. */
  comment: string;
  /** The dependent_columns is the list of dependent columns of a view. */
  dependentColumns: DependentColumn[];
}

/** DependentColumn is the metadata for dependent columns. */
export interface DependentColumn {
  /** The schema is the schema of a reference column. */
  schema: string;
  /** The table is the table of a reference column. */
  table: string;
  /** The column is the name of a reference column. */
  column: string;
}

/** FunctionMetadata is the metadata for functions. */
export interface FunctionMetadata {
  /** The name is the name of a view. */
  name: string;
  /** The definition is the definition of a view. */
  definition: string;
}

/** IndexMetadata is the metadata for indexes. */
export interface IndexMetadata {
  /** The name is the name of an index. */
  name: string;
  /**
   * The expressions are the ordered columns or expressions of an index.
   * This could refer to a column or an expression.
   */
  expressions: string[];
  /** The type is the type of an index. */
  type: string;
  /** The unique is whether the index is unique. */
  unique: boolean;
  /** The primary is whether the index is a primary key index. */
  primary: boolean;
  /** The visible is whether the index is visible. */
  visible: boolean;
  /** The comment is the comment of an index. */
  comment: string;
}

/** ExtensionMetadata is the metadata for extensions. */
export interface ExtensionMetadata {
  /** The name is the name of an extension. */
  name: string;
  /** The schema is the extension that is installed to. But the extension usage is not limited to the schema. */
  schema: string;
  /** The version is the version of an extension. */
  version: string;
  /** The description is the description of an extension. */
  description: string;
}

/** ForeignKeyMetadata is the metadata for foreign keys. */
export interface ForeignKeyMetadata {
  /** The name is the name of a foreign key. */
  name: string;
  /** The columns are the ordered referencing columns of a foreign key. */
  columns: string[];
  /**
   * The referenced_schema is the referenced schema name of a foreign key.
   * It is an empty string for databases without such concept such as MySQL.
   */
  referencedSchema: string;
  /** The referenced_table is the referenced table name of a foreign key. */
  referencedTable: string;
  /** The referenced_columns are the ordered referenced columns of a foreign key. */
  referencedColumns: string[];
  /** The on_delete is the on delete action of a foreign key. */
  onDelete: string;
  /** The on_update is the on update action of a foreign key. */
  onUpdate: string;
  /**
   * The match_type is the match type of a foreign key.
   * The match_type is the PostgreSQL specific field.
   * It's empty string for other databases.
   */
  matchType: string;
}

/** InstanceRoleMetadata is the message for instance role. */
export interface InstanceRoleMetadata {
  /** The role name. It's unique within the instance. */
  name: string;
  /** The grant display string on the instance. It's generated by database engine. */
  grant: string;
}

export interface Secrets {
  /** The list of secrets. */
  items: SecretItem[];
}

export interface SecretItem {
  /** The name is the name of the secret. */
  name: string;
  /** The value is the value of the secret. */
  value: string;
  /** The description is the description of the secret. */
  description: string;
}

function createBaseDatabaseMetadata(): DatabaseMetadata {
  return { name: "", schemas: [], characterSet: "", collation: "", extensions: [] };
}

export const DatabaseMetadata = {
  encode(message: DatabaseMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    for (const v of message.schemas) {
      SchemaMetadata.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.characterSet !== "") {
      writer.uint32(26).string(message.characterSet);
    }
    if (message.collation !== "") {
      writer.uint32(34).string(message.collation);
    }
    for (const v of message.extensions) {
      ExtensionMetadata.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DatabaseMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDatabaseMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.schemas.push(SchemaMetadata.decode(reader, reader.uint32()));
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.characterSet = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.collation = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.extensions.push(ExtensionMetadata.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DatabaseMetadata {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      schemas: Array.isArray(object?.schemas) ? object.schemas.map((e: any) => SchemaMetadata.fromJSON(e)) : [],
      characterSet: isSet(object.characterSet) ? String(object.characterSet) : "",
      collation: isSet(object.collation) ? String(object.collation) : "",
      extensions: Array.isArray(object?.extensions)
        ? object.extensions.map((e: any) => ExtensionMetadata.fromJSON(e))
        : [],
    };
  },

  toJSON(message: DatabaseMetadata): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    if (message.schemas) {
      obj.schemas = message.schemas.map((e) => e ? SchemaMetadata.toJSON(e) : undefined);
    } else {
      obj.schemas = [];
    }
    message.characterSet !== undefined && (obj.characterSet = message.characterSet);
    message.collation !== undefined && (obj.collation = message.collation);
    if (message.extensions) {
      obj.extensions = message.extensions.map((e) => e ? ExtensionMetadata.toJSON(e) : undefined);
    } else {
      obj.extensions = [];
    }
    return obj;
  },

  create(base?: DeepPartial<DatabaseMetadata>): DatabaseMetadata {
    return DatabaseMetadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DatabaseMetadata>): DatabaseMetadata {
    const message = createBaseDatabaseMetadata();
    message.name = object.name ?? "";
    message.schemas = object.schemas?.map((e) => SchemaMetadata.fromPartial(e)) || [];
    message.characterSet = object.characterSet ?? "";
    message.collation = object.collation ?? "";
    message.extensions = object.extensions?.map((e) => ExtensionMetadata.fromPartial(e)) || [];
    return message;
  },
};

function createBaseSchemaMetadata(): SchemaMetadata {
  return { name: "", tables: [], views: [], functions: [] };
}

export const SchemaMetadata = {
  encode(message: SchemaMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    for (const v of message.tables) {
      TableMetadata.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.views) {
      ViewMetadata.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.functions) {
      FunctionMetadata.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SchemaMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSchemaMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.tables.push(TableMetadata.decode(reader, reader.uint32()));
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.views.push(ViewMetadata.decode(reader, reader.uint32()));
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.functions.push(FunctionMetadata.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SchemaMetadata {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      tables: Array.isArray(object?.tables) ? object.tables.map((e: any) => TableMetadata.fromJSON(e)) : [],
      views: Array.isArray(object?.views) ? object.views.map((e: any) => ViewMetadata.fromJSON(e)) : [],
      functions: Array.isArray(object?.functions) ? object.functions.map((e: any) => FunctionMetadata.fromJSON(e)) : [],
    };
  },

  toJSON(message: SchemaMetadata): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    if (message.tables) {
      obj.tables = message.tables.map((e) => e ? TableMetadata.toJSON(e) : undefined);
    } else {
      obj.tables = [];
    }
    if (message.views) {
      obj.views = message.views.map((e) => e ? ViewMetadata.toJSON(e) : undefined);
    } else {
      obj.views = [];
    }
    if (message.functions) {
      obj.functions = message.functions.map((e) => e ? FunctionMetadata.toJSON(e) : undefined);
    } else {
      obj.functions = [];
    }
    return obj;
  },

  create(base?: DeepPartial<SchemaMetadata>): SchemaMetadata {
    return SchemaMetadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SchemaMetadata>): SchemaMetadata {
    const message = createBaseSchemaMetadata();
    message.name = object.name ?? "";
    message.tables = object.tables?.map((e) => TableMetadata.fromPartial(e)) || [];
    message.views = object.views?.map((e) => ViewMetadata.fromPartial(e)) || [];
    message.functions = object.functions?.map((e) => FunctionMetadata.fromPartial(e)) || [];
    return message;
  },
};

function createBaseTableMetadata(): TableMetadata {
  return {
    name: "",
    columns: [],
    indexes: [],
    engine: "",
    collation: "",
    rowCount: 0,
    dataSize: 0,
    indexSize: 0,
    dataFree: 0,
    createOptions: "",
    comment: "",
    foreignKeys: [],
  };
}

export const TableMetadata = {
  encode(message: TableMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    for (const v of message.columns) {
      ColumnMetadata.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.indexes) {
      IndexMetadata.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.engine !== "") {
      writer.uint32(34).string(message.engine);
    }
    if (message.collation !== "") {
      writer.uint32(42).string(message.collation);
    }
    if (message.rowCount !== 0) {
      writer.uint32(48).int64(message.rowCount);
    }
    if (message.dataSize !== 0) {
      writer.uint32(56).int64(message.dataSize);
    }
    if (message.indexSize !== 0) {
      writer.uint32(64).int64(message.indexSize);
    }
    if (message.dataFree !== 0) {
      writer.uint32(72).int64(message.dataFree);
    }
    if (message.createOptions !== "") {
      writer.uint32(82).string(message.createOptions);
    }
    if (message.comment !== "") {
      writer.uint32(90).string(message.comment);
    }
    for (const v of message.foreignKeys) {
      ForeignKeyMetadata.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TableMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTableMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.columns.push(ColumnMetadata.decode(reader, reader.uint32()));
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.indexes.push(IndexMetadata.decode(reader, reader.uint32()));
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.engine = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.collation = reader.string();
          continue;
        case 6:
          if (tag != 48) {
            break;
          }

          message.rowCount = longToNumber(reader.int64() as Long);
          continue;
        case 7:
          if (tag != 56) {
            break;
          }

          message.dataSize = longToNumber(reader.int64() as Long);
          continue;
        case 8:
          if (tag != 64) {
            break;
          }

          message.indexSize = longToNumber(reader.int64() as Long);
          continue;
        case 9:
          if (tag != 72) {
            break;
          }

          message.dataFree = longToNumber(reader.int64() as Long);
          continue;
        case 10:
          if (tag != 82) {
            break;
          }

          message.createOptions = reader.string();
          continue;
        case 11:
          if (tag != 90) {
            break;
          }

          message.comment = reader.string();
          continue;
        case 12:
          if (tag != 98) {
            break;
          }

          message.foreignKeys.push(ForeignKeyMetadata.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TableMetadata {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      columns: Array.isArray(object?.columns) ? object.columns.map((e: any) => ColumnMetadata.fromJSON(e)) : [],
      indexes: Array.isArray(object?.indexes) ? object.indexes.map((e: any) => IndexMetadata.fromJSON(e)) : [],
      engine: isSet(object.engine) ? String(object.engine) : "",
      collation: isSet(object.collation) ? String(object.collation) : "",
      rowCount: isSet(object.rowCount) ? Number(object.rowCount) : 0,
      dataSize: isSet(object.dataSize) ? Number(object.dataSize) : 0,
      indexSize: isSet(object.indexSize) ? Number(object.indexSize) : 0,
      dataFree: isSet(object.dataFree) ? Number(object.dataFree) : 0,
      createOptions: isSet(object.createOptions) ? String(object.createOptions) : "",
      comment: isSet(object.comment) ? String(object.comment) : "",
      foreignKeys: Array.isArray(object?.foreignKeys)
        ? object.foreignKeys.map((e: any) => ForeignKeyMetadata.fromJSON(e))
        : [],
    };
  },

  toJSON(message: TableMetadata): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    if (message.columns) {
      obj.columns = message.columns.map((e) => e ? ColumnMetadata.toJSON(e) : undefined);
    } else {
      obj.columns = [];
    }
    if (message.indexes) {
      obj.indexes = message.indexes.map((e) => e ? IndexMetadata.toJSON(e) : undefined);
    } else {
      obj.indexes = [];
    }
    message.engine !== undefined && (obj.engine = message.engine);
    message.collation !== undefined && (obj.collation = message.collation);
    message.rowCount !== undefined && (obj.rowCount = Math.round(message.rowCount));
    message.dataSize !== undefined && (obj.dataSize = Math.round(message.dataSize));
    message.indexSize !== undefined && (obj.indexSize = Math.round(message.indexSize));
    message.dataFree !== undefined && (obj.dataFree = Math.round(message.dataFree));
    message.createOptions !== undefined && (obj.createOptions = message.createOptions);
    message.comment !== undefined && (obj.comment = message.comment);
    if (message.foreignKeys) {
      obj.foreignKeys = message.foreignKeys.map((e) => e ? ForeignKeyMetadata.toJSON(e) : undefined);
    } else {
      obj.foreignKeys = [];
    }
    return obj;
  },

  create(base?: DeepPartial<TableMetadata>): TableMetadata {
    return TableMetadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TableMetadata>): TableMetadata {
    const message = createBaseTableMetadata();
    message.name = object.name ?? "";
    message.columns = object.columns?.map((e) => ColumnMetadata.fromPartial(e)) || [];
    message.indexes = object.indexes?.map((e) => IndexMetadata.fromPartial(e)) || [];
    message.engine = object.engine ?? "";
    message.collation = object.collation ?? "";
    message.rowCount = object.rowCount ?? 0;
    message.dataSize = object.dataSize ?? 0;
    message.indexSize = object.indexSize ?? 0;
    message.dataFree = object.dataFree ?? 0;
    message.createOptions = object.createOptions ?? "";
    message.comment = object.comment ?? "";
    message.foreignKeys = object.foreignKeys?.map((e) => ForeignKeyMetadata.fromPartial(e)) || [];
    return message;
  },
};

function createBaseColumnMetadata(): ColumnMetadata {
  return {
    name: "",
    position: 0,
    default: undefined,
    nullable: false,
    type: "",
    characterSet: "",
    collation: "",
    comment: "",
  };
}

export const ColumnMetadata = {
  encode(message: ColumnMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.position !== 0) {
      writer.uint32(16).int32(message.position);
    }
    if (message.default !== undefined) {
      StringValue.encode({ value: message.default! }, writer.uint32(26).fork()).ldelim();
    }
    if (message.nullable === true) {
      writer.uint32(32).bool(message.nullable);
    }
    if (message.type !== "") {
      writer.uint32(42).string(message.type);
    }
    if (message.characterSet !== "") {
      writer.uint32(50).string(message.characterSet);
    }
    if (message.collation !== "") {
      writer.uint32(58).string(message.collation);
    }
    if (message.comment !== "") {
      writer.uint32(66).string(message.comment);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ColumnMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseColumnMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.position = reader.int32();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.default = StringValue.decode(reader, reader.uint32()).value;
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.nullable = reader.bool();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.type = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.characterSet = reader.string();
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.collation = reader.string();
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.comment = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ColumnMetadata {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      position: isSet(object.position) ? Number(object.position) : 0,
      default: isSet(object.default) ? String(object.default) : undefined,
      nullable: isSet(object.nullable) ? Boolean(object.nullable) : false,
      type: isSet(object.type) ? String(object.type) : "",
      characterSet: isSet(object.characterSet) ? String(object.characterSet) : "",
      collation: isSet(object.collation) ? String(object.collation) : "",
      comment: isSet(object.comment) ? String(object.comment) : "",
    };
  },

  toJSON(message: ColumnMetadata): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.position !== undefined && (obj.position = Math.round(message.position));
    message.default !== undefined && (obj.default = message.default);
    message.nullable !== undefined && (obj.nullable = message.nullable);
    message.type !== undefined && (obj.type = message.type);
    message.characterSet !== undefined && (obj.characterSet = message.characterSet);
    message.collation !== undefined && (obj.collation = message.collation);
    message.comment !== undefined && (obj.comment = message.comment);
    return obj;
  },

  create(base?: DeepPartial<ColumnMetadata>): ColumnMetadata {
    return ColumnMetadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ColumnMetadata>): ColumnMetadata {
    const message = createBaseColumnMetadata();
    message.name = object.name ?? "";
    message.position = object.position ?? 0;
    message.default = object.default ?? undefined;
    message.nullable = object.nullable ?? false;
    message.type = object.type ?? "";
    message.characterSet = object.characterSet ?? "";
    message.collation = object.collation ?? "";
    message.comment = object.comment ?? "";
    return message;
  },
};

function createBaseViewMetadata(): ViewMetadata {
  return { name: "", definition: "", comment: "", dependentColumns: [] };
}

export const ViewMetadata = {
  encode(message: ViewMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.definition !== "") {
      writer.uint32(18).string(message.definition);
    }
    if (message.comment !== "") {
      writer.uint32(26).string(message.comment);
    }
    for (const v of message.dependentColumns) {
      DependentColumn.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ViewMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseViewMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.definition = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.comment = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.dependentColumns.push(DependentColumn.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ViewMetadata {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      definition: isSet(object.definition) ? String(object.definition) : "",
      comment: isSet(object.comment) ? String(object.comment) : "",
      dependentColumns: Array.isArray(object?.dependentColumns)
        ? object.dependentColumns.map((e: any) => DependentColumn.fromJSON(e))
        : [],
    };
  },

  toJSON(message: ViewMetadata): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.definition !== undefined && (obj.definition = message.definition);
    message.comment !== undefined && (obj.comment = message.comment);
    if (message.dependentColumns) {
      obj.dependentColumns = message.dependentColumns.map((e) => e ? DependentColumn.toJSON(e) : undefined);
    } else {
      obj.dependentColumns = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ViewMetadata>): ViewMetadata {
    return ViewMetadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ViewMetadata>): ViewMetadata {
    const message = createBaseViewMetadata();
    message.name = object.name ?? "";
    message.definition = object.definition ?? "";
    message.comment = object.comment ?? "";
    message.dependentColumns = object.dependentColumns?.map((e) => DependentColumn.fromPartial(e)) || [];
    return message;
  },
};

function createBaseDependentColumn(): DependentColumn {
  return { schema: "", table: "", column: "" };
}

export const DependentColumn = {
  encode(message: DependentColumn, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.schema !== "") {
      writer.uint32(10).string(message.schema);
    }
    if (message.table !== "") {
      writer.uint32(18).string(message.table);
    }
    if (message.column !== "") {
      writer.uint32(26).string(message.column);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DependentColumn {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDependentColumn();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.schema = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.table = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.column = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DependentColumn {
    return {
      schema: isSet(object.schema) ? String(object.schema) : "",
      table: isSet(object.table) ? String(object.table) : "",
      column: isSet(object.column) ? String(object.column) : "",
    };
  },

  toJSON(message: DependentColumn): unknown {
    const obj: any = {};
    message.schema !== undefined && (obj.schema = message.schema);
    message.table !== undefined && (obj.table = message.table);
    message.column !== undefined && (obj.column = message.column);
    return obj;
  },

  create(base?: DeepPartial<DependentColumn>): DependentColumn {
    return DependentColumn.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DependentColumn>): DependentColumn {
    const message = createBaseDependentColumn();
    message.schema = object.schema ?? "";
    message.table = object.table ?? "";
    message.column = object.column ?? "";
    return message;
  },
};

function createBaseFunctionMetadata(): FunctionMetadata {
  return { name: "", definition: "" };
}

export const FunctionMetadata = {
  encode(message: FunctionMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.definition !== "") {
      writer.uint32(18).string(message.definition);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FunctionMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFunctionMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.definition = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FunctionMetadata {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      definition: isSet(object.definition) ? String(object.definition) : "",
    };
  },

  toJSON(message: FunctionMetadata): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.definition !== undefined && (obj.definition = message.definition);
    return obj;
  },

  create(base?: DeepPartial<FunctionMetadata>): FunctionMetadata {
    return FunctionMetadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<FunctionMetadata>): FunctionMetadata {
    const message = createBaseFunctionMetadata();
    message.name = object.name ?? "";
    message.definition = object.definition ?? "";
    return message;
  },
};

function createBaseIndexMetadata(): IndexMetadata {
  return { name: "", expressions: [], type: "", unique: false, primary: false, visible: false, comment: "" };
}

export const IndexMetadata = {
  encode(message: IndexMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    for (const v of message.expressions) {
      writer.uint32(18).string(v!);
    }
    if (message.type !== "") {
      writer.uint32(26).string(message.type);
    }
    if (message.unique === true) {
      writer.uint32(32).bool(message.unique);
    }
    if (message.primary === true) {
      writer.uint32(40).bool(message.primary);
    }
    if (message.visible === true) {
      writer.uint32(48).bool(message.visible);
    }
    if (message.comment !== "") {
      writer.uint32(58).string(message.comment);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IndexMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIndexMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.expressions.push(reader.string());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.type = reader.string();
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.unique = reader.bool();
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.primary = reader.bool();
          continue;
        case 6:
          if (tag != 48) {
            break;
          }

          message.visible = reader.bool();
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.comment = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IndexMetadata {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      expressions: Array.isArray(object?.expressions) ? object.expressions.map((e: any) => String(e)) : [],
      type: isSet(object.type) ? String(object.type) : "",
      unique: isSet(object.unique) ? Boolean(object.unique) : false,
      primary: isSet(object.primary) ? Boolean(object.primary) : false,
      visible: isSet(object.visible) ? Boolean(object.visible) : false,
      comment: isSet(object.comment) ? String(object.comment) : "",
    };
  },

  toJSON(message: IndexMetadata): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    if (message.expressions) {
      obj.expressions = message.expressions.map((e) => e);
    } else {
      obj.expressions = [];
    }
    message.type !== undefined && (obj.type = message.type);
    message.unique !== undefined && (obj.unique = message.unique);
    message.primary !== undefined && (obj.primary = message.primary);
    message.visible !== undefined && (obj.visible = message.visible);
    message.comment !== undefined && (obj.comment = message.comment);
    return obj;
  },

  create(base?: DeepPartial<IndexMetadata>): IndexMetadata {
    return IndexMetadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IndexMetadata>): IndexMetadata {
    const message = createBaseIndexMetadata();
    message.name = object.name ?? "";
    message.expressions = object.expressions?.map((e) => e) || [];
    message.type = object.type ?? "";
    message.unique = object.unique ?? false;
    message.primary = object.primary ?? false;
    message.visible = object.visible ?? false;
    message.comment = object.comment ?? "";
    return message;
  },
};

function createBaseExtensionMetadata(): ExtensionMetadata {
  return { name: "", schema: "", version: "", description: "" };
}

export const ExtensionMetadata = {
  encode(message: ExtensionMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.schema !== "") {
      writer.uint32(18).string(message.schema);
    }
    if (message.version !== "") {
      writer.uint32(26).string(message.version);
    }
    if (message.description !== "") {
      writer.uint32(34).string(message.description);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExtensionMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExtensionMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.schema = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.version = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.description = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExtensionMetadata {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      schema: isSet(object.schema) ? String(object.schema) : "",
      version: isSet(object.version) ? String(object.version) : "",
      description: isSet(object.description) ? String(object.description) : "",
    };
  },

  toJSON(message: ExtensionMetadata): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.schema !== undefined && (obj.schema = message.schema);
    message.version !== undefined && (obj.version = message.version);
    message.description !== undefined && (obj.description = message.description);
    return obj;
  },

  create(base?: DeepPartial<ExtensionMetadata>): ExtensionMetadata {
    return ExtensionMetadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExtensionMetadata>): ExtensionMetadata {
    const message = createBaseExtensionMetadata();
    message.name = object.name ?? "";
    message.schema = object.schema ?? "";
    message.version = object.version ?? "";
    message.description = object.description ?? "";
    return message;
  },
};

function createBaseForeignKeyMetadata(): ForeignKeyMetadata {
  return {
    name: "",
    columns: [],
    referencedSchema: "",
    referencedTable: "",
    referencedColumns: [],
    onDelete: "",
    onUpdate: "",
    matchType: "",
  };
}

export const ForeignKeyMetadata = {
  encode(message: ForeignKeyMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    for (const v of message.columns) {
      writer.uint32(18).string(v!);
    }
    if (message.referencedSchema !== "") {
      writer.uint32(26).string(message.referencedSchema);
    }
    if (message.referencedTable !== "") {
      writer.uint32(34).string(message.referencedTable);
    }
    for (const v of message.referencedColumns) {
      writer.uint32(42).string(v!);
    }
    if (message.onDelete !== "") {
      writer.uint32(50).string(message.onDelete);
    }
    if (message.onUpdate !== "") {
      writer.uint32(58).string(message.onUpdate);
    }
    if (message.matchType !== "") {
      writer.uint32(66).string(message.matchType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ForeignKeyMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseForeignKeyMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.columns.push(reader.string());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.referencedSchema = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.referencedTable = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.referencedColumns.push(reader.string());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.onDelete = reader.string();
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.onUpdate = reader.string();
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.matchType = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ForeignKeyMetadata {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      columns: Array.isArray(object?.columns) ? object.columns.map((e: any) => String(e)) : [],
      referencedSchema: isSet(object.referencedSchema) ? String(object.referencedSchema) : "",
      referencedTable: isSet(object.referencedTable) ? String(object.referencedTable) : "",
      referencedColumns: Array.isArray(object?.referencedColumns)
        ? object.referencedColumns.map((e: any) => String(e))
        : [],
      onDelete: isSet(object.onDelete) ? String(object.onDelete) : "",
      onUpdate: isSet(object.onUpdate) ? String(object.onUpdate) : "",
      matchType: isSet(object.matchType) ? String(object.matchType) : "",
    };
  },

  toJSON(message: ForeignKeyMetadata): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    if (message.columns) {
      obj.columns = message.columns.map((e) => e);
    } else {
      obj.columns = [];
    }
    message.referencedSchema !== undefined && (obj.referencedSchema = message.referencedSchema);
    message.referencedTable !== undefined && (obj.referencedTable = message.referencedTable);
    if (message.referencedColumns) {
      obj.referencedColumns = message.referencedColumns.map((e) => e);
    } else {
      obj.referencedColumns = [];
    }
    message.onDelete !== undefined && (obj.onDelete = message.onDelete);
    message.onUpdate !== undefined && (obj.onUpdate = message.onUpdate);
    message.matchType !== undefined && (obj.matchType = message.matchType);
    return obj;
  },

  create(base?: DeepPartial<ForeignKeyMetadata>): ForeignKeyMetadata {
    return ForeignKeyMetadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ForeignKeyMetadata>): ForeignKeyMetadata {
    const message = createBaseForeignKeyMetadata();
    message.name = object.name ?? "";
    message.columns = object.columns?.map((e) => e) || [];
    message.referencedSchema = object.referencedSchema ?? "";
    message.referencedTable = object.referencedTable ?? "";
    message.referencedColumns = object.referencedColumns?.map((e) => e) || [];
    message.onDelete = object.onDelete ?? "";
    message.onUpdate = object.onUpdate ?? "";
    message.matchType = object.matchType ?? "";
    return message;
  },
};

function createBaseInstanceRoleMetadata(): InstanceRoleMetadata {
  return { name: "", grant: "" };
}

export const InstanceRoleMetadata = {
  encode(message: InstanceRoleMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.grant !== "") {
      writer.uint32(58).string(message.grant);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InstanceRoleMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstanceRoleMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          if (tag != 18) {
            break;
          }

          message.name = reader.string();
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.grant = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): InstanceRoleMetadata {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      grant: isSet(object.grant) ? String(object.grant) : "",
    };
  },

  toJSON(message: InstanceRoleMetadata): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.grant !== undefined && (obj.grant = message.grant);
    return obj;
  },

  create(base?: DeepPartial<InstanceRoleMetadata>): InstanceRoleMetadata {
    return InstanceRoleMetadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InstanceRoleMetadata>): InstanceRoleMetadata {
    const message = createBaseInstanceRoleMetadata();
    message.name = object.name ?? "";
    message.grant = object.grant ?? "";
    return message;
  },
};

function createBaseSecrets(): Secrets {
  return { items: [] };
}

export const Secrets = {
  encode(message: Secrets, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.items) {
      SecretItem.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Secrets {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecrets();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.items.push(SecretItem.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Secrets {
    return { items: Array.isArray(object?.items) ? object.items.map((e: any) => SecretItem.fromJSON(e)) : [] };
  },

  toJSON(message: Secrets): unknown {
    const obj: any = {};
    if (message.items) {
      obj.items = message.items.map((e) => e ? SecretItem.toJSON(e) : undefined);
    } else {
      obj.items = [];
    }
    return obj;
  },

  create(base?: DeepPartial<Secrets>): Secrets {
    return Secrets.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Secrets>): Secrets {
    const message = createBaseSecrets();
    message.items = object.items?.map((e) => SecretItem.fromPartial(e)) || [];
    return message;
  },
};

function createBaseSecretItem(): SecretItem {
  return { name: "", value: "", description: "" };
}

export const SecretItem = {
  encode(message: SecretItem, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    if (message.description !== "") {
      writer.uint32(26).string(message.description);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecretItem {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecretItem();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.description = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SecretItem {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      value: isSet(object.value) ? String(object.value) : "",
      description: isSet(object.description) ? String(object.description) : "",
    };
  },

  toJSON(message: SecretItem): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.value !== undefined && (obj.value = message.value);
    message.description !== undefined && (obj.description = message.description);
    return obj;
  },

  create(base?: DeepPartial<SecretItem>): SecretItem {
    return SecretItem.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SecretItem>): SecretItem {
    const message = createBaseSecretItem();
    message.name = object.name ?? "";
    message.value = object.value ?? "";
    message.description = object.description ?? "";
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

// If you get a compile-error about 'Constructor<Long> and ... have no overlap',
// add '--ts_proto_opt=esModuleInterop=true' as a flag when calling 'protoc'.
if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
